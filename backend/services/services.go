package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"evoting_server/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/hyperledger/fabric-gateway/pkg/client"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type VotingService struct {
	Contract          *client.Contract
	VotersColl        *mongo.Collection
	CandidatesColl    *mongo.Collection
	MongoDB           *mongo.Database
	S3Client          *s3.Client
	RekognitionClient *rekognition.Client
	S3Bucket          string
}

func NewVotingService(
	voters *mongo.Collection,
	candidates *mongo.Collection,
	contract *client.Contract,
	s3Client *s3.Client,
	rekogClient *rekognition.Client,
	s3Bucket string,
) *VotingService {
	return &VotingService{
		VotersColl:        voters,
		CandidatesColl:    candidates,
		Contract:          contract,
		S3Client:          s3Client,
		RekognitionClient: rekogClient,
		S3Bucket:          s3Bucket,
	}
}

func (s *VotingService) GetCandidatesByVoterService(ctx context.Context, voterID string) ([]models.Candidate, error) {

	voterCollection := s.VotersColl

	var voter models.Voter
	err := voterCollection.FindOne(ctx, bson.M{"_id": voterID}).Decode(&voter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("voter not found")
		}
		return nil, fmt.Errorf("failed to get voter constituency: %v", err)
	}

	candidateCollection := s.CandidatesColl
	cursor, err := candidateCollection.Find(ctx, bson.M{"constituency": voter.Constituency})
	if err != nil {
		return nil, fmt.Errorf("failed to get candidates: %v", err)
	}
	defer cursor.Close(ctx)

	var candidates []models.Candidate
	if err := cursor.All(ctx, &candidates); err != nil {
		return nil, fmt.Errorf("failed to decode candidates: %v", err)
	}

	return candidates, nil
}

func (s *VotingService) RegisterVoter(ctx context.Context, voterID, name, constituency string, age int, imageBytes []byte) error {

	// Reject underage voters
	if age < 18 {
		return errors.New("voter must be at least 18 years old")
	}

	// Check if voter exists
	var existing models.Voter
	err := s.VotersColl.FindOne(ctx, bson.M{"_id": voterID}).Decode(&existing)
	if err == nil {
		return errors.New("voter already registered")
	}
	if err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking existing voter: %v", err)
	}

	s3Key := fmt.Sprintf("%s/%s.jpg", constituency, voterID)
	err = s.UploadImageToS3(s3Key, imageBytes)
	if err != nil {
		return fmt.Errorf("failed to upload face image to S3: %v", err)
	}

	// 1. Register voter on Fabric blockchain
	_, err = s.Contract.SubmitTransaction("RegisterVoter", voterID, constituency)
	if err != nil {
		return fmt.Errorf("chaincode error: %v", err)
	}

	// 2. Store voter in MongoDB
	voter := models.Voter{
		ID:           voterID,
		Name:         name,
		Constituency: constituency,
		Age:          age,
		HasVoted:     false,
	}

	_, err = s.VotersColl.InsertOne(ctx, voter)
	if err != nil {
		return fmt.Errorf("mongo insert error: %v", err)
	}

	return nil
}

func (s *VotingService) CastVote(voterID string, candidateID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var voter models.Voter
	err := s.VotersColl.FindOne(ctx, bson.M{"_id": voterID}).Decode(&voter)
	if err != nil {
		return errors.New("voter not found")
	}
	if voter.HasVoted {
		return errors.New("voter has already voted")
	}

	_, err = s.Contract.SubmitTransaction("CastVote", voterID, candidateID)
	if err != nil {
		return fmt.Errorf("chaincode error: %v", err)
	}

	_, err = s.VotersColl.UpdateOne(ctx, bson.M{"_id": voterID}, bson.M{"$set": bson.M{"hasVoted": true}})
	return err
}

func (s *VotingService) GetResults() (map[string]int, error) {
	evaluateResult, err := s.Contract.EvaluateTransaction("GetResults")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate chaincode function: %v", err)
	}

	var results map[string]int
	err = json.Unmarshal(evaluateResult, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chaincode result: %v", err)
	}

	return results, nil
}

func (s *VotingService) GetResultsByVoterService(voterID string) ([]models.Candidate, error) {
	// Step 1: Get candidates from chaincode
	candidatesBytes, err := s.Contract.EvaluateTransaction("GetResultByVoterID", voterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidates from chaincode: %v", err)
	}

	var candidates []models.Candidate
	if err := json.Unmarshal(candidatesBytes, &candidates); err != nil {
		return nil, fmt.Errorf("failed to parse candidates: %v", err)
	}

	// Step 2: Build list of IDs to query MongoDB
	ids := make([]string, len(candidates))
	for i, c := range candidates {
		ids[i] = c.ID
	}

	// Step 3: Query MongoDB for names
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := s.CandidatesColl.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch candidate names from MongoDB: %v", err)
	}
	defer cursor.Close(context.TODO())

	nameMap := make(map[string]string)
	for cursor.Next(context.TODO()) {
		var mongoCandidate models.Candidate
		if err := cursor.Decode(&mongoCandidate); err != nil {
			continue // skip invalid
		}
		nameMap[mongoCandidate.ID] = mongoCandidate.Name
	}

	// Step 4: Merge name into chaincode result
	for i := range candidates {
		if name, exists := nameMap[candidates[i].ID]; exists {
			candidates[i].Name = name
		}
	}

	return candidates, nil
}

// VerifyVoter checks if a voterID exists and face match is true
func (s *VotingService) VerifyVoter(voterID string, imageBytes []byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var voter models.Voter
	err := s.VotersColl.FindOne(ctx, bson.M{"_id": voterID}).Decode(&voter)
	if err != nil {
		return false, errors.New("voter not found")
	}

	// Assuming voter's image is stored in S3 as "{Constituency}/<voterID>.jpg"
	s3Key := fmt.Sprintf("%s/%s.jpg", voter.Constituency, voterID)

	// Compare using Rekognition
	match, err := s.CompareFacesWithS3(s3Key, imageBytes)
	if err != nil {
		return false, fmt.Errorf("face comparison failed: %v", err)
	}

	if !match {
		return false, errors.New("face mismatch")
	}

	return true, nil
}

// Functions for Cloud Service : AWS
func (s *VotingService) UploadImageToS3(key string, imageBytes []byte) error {
	_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.S3Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(imageBytes),
		ContentType: aws.String("image/jpeg"),
	})
	return err
}

func (s *VotingService) CompareFacesWithS3(sourceKey string, targetImageBytes []byte) (bool, error) {
	input := &rekognition.CompareFacesInput{
		SourceImage: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(s.S3Bucket),
				Name:   aws.String(sourceKey),
			},
		},
		TargetImage: &types.Image{
			Bytes: targetImageBytes,
		},
		SimilarityThreshold: aws.Float32(90.0),
	}

	output, err := s.RekognitionClient.CompareFaces(context.TODO(), input)
	if err != nil {
		return false, err
	}

	for _, match := range output.FaceMatches {
		if match.Similarity != nil && *match.Similarity >= 90.0 {
			return true, nil
		}
	}

	return false, nil
}
