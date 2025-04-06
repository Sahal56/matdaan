package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

/*
operations:
    Register Voter
	Cast Vote
	GetResultsofAllCandidates
	GetResultofOneCandidate


Scalability Criteria | By channel based on Constituency/Region
*/

type SmartContract struct {
	contractapi.Contract
}

type Voter struct {
	ID           string `json:"ID"`
	Constituency string `json:"Constituency"`
	HasVoted     bool   `json:"HasVoted"`
}

type Candidate struct {
	ID           string `bson:"_id,omitempty" json:"ID"`
	Constituency string `json:"Constituency"`
	Votes        int    `json:"Votes"`
	HasVoted     bool   `json:"HasVoted"` // new logic: candidates can also vote
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://host.docker.internal:27017"))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("evotingDB").Collection("candidates")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return fmt.Errorf("failed to fetch candidates from DB: %v", err)
	}
	defer cursor.Close(context.TODO())

	var candidates []Candidate
	for cursor.Next(context.TODO()) {
		var candidate Candidate
		err := cursor.Decode(&candidate)
		if err != nil {
			return fmt.Errorf("failed to decode candidate: %v", err)
		}
		candidates = append(candidates, candidate)
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("error iterating cursor: %v", err)
	}

	numCandidates := len(candidates)

	// Store candidates dynamically
	for _, candidate := range candidates {
		candidateBytes, _ := json.Marshal(candidate)
		err := ctx.GetStub().PutState(candidate.ID, candidateBytes)
		if err != nil {
			return fmt.Errorf("failed to put state for %s: %v", candidate.ID, err)
		}
	}

	// Store candidate count
	err = ctx.GetStub().PutState("CANDIDATE_COUNT", []byte(strconv.Itoa(numCandidates)))
	if err != nil {
		return fmt.Errorf("failed to store candidate count: %v", err)
	}

	return nil
}

func (s *SmartContract) RegisterVoter(ctx contractapi.TransactionContextInterface, voterID string, constituency string) error {

	voter := Voter{
		ID:           voterID,
		Constituency: constituency,
		HasVoted:     false,
	}

	voterBytes, _ := json.Marshal(voter)
	err := ctx.GetStub().PutState(voterID, voterBytes)
	if err != nil {
		return fmt.Errorf("failed to register voter %s: %v", voterID, err)
	}
	return nil
}

func (s *SmartContract) CastVote(ctx contractapi.TransactionContextInterface, voterID string, candidateID string) error {
	voterBytes, err := ctx.GetStub().GetState(voterID)
	if err != nil {
		return fmt.Errorf("failed to read voter state: %v", err)
	}
	if voterBytes == nil {
		return fmt.Errorf("voter %s does not exist", voterID)
	}

	var voter Voter
	json.Unmarshal(voterBytes, &voter)

	if voter.HasVoted {
		return fmt.Errorf("voter %s has already voted", voterID)
	}

	candidateBytes, err := ctx.GetStub().GetState(candidateID)
	if err != nil {
		return fmt.Errorf("failed to read candidate state: %v", err)
	}
	if candidateBytes == nil {
		return fmt.Errorf("candidate %s does not exist", candidateID)
	}

	var candidate Candidate
	json.Unmarshal(candidateBytes, &candidate)

	// new code ? research
	// although it is verified at front end
	// But it is ensured here. Voter can only vote for his/her constituency
	if voter.Constituency != candidate.Constituency {
		return fmt.Errorf("voter:%s is not allowed to vote for candidate of another constituency", voterID)
	}

	candidate.Votes = candidate.Votes + 1
	candidateBytes, _ = json.Marshal(candidate)
	err = ctx.GetStub().PutState(candidateID, candidateBytes)
	if err != nil {
		return fmt.Errorf("failed to update candidate votes: %v", err)
	}

	voter.HasVoted = true
	voterBytes, _ = json.Marshal(voter)
	err = ctx.GetStub().PutState(voterID, voterBytes)
	if err != nil {
		return fmt.Errorf("failed to update voter status: %v", err)
	}

	return nil
}

func (s *SmartContract) CastVoteByCandidate(ctx contractapi.TransactionContextInterface, sending_candidateID string, receiving_candidateID string) error {
	// Candidate can also vote.
	sending_candidateBytes, err := ctx.GetStub().GetState(sending_candidateID)
	if err != nil {
		return fmt.Errorf("failed to read candidate state: %v", err)
	}
	if sending_candidateBytes == nil {
		return fmt.Errorf("candidate %s does not exist", sending_candidateID)
	}
	var sending_candidate Candidate
	json.Unmarshal(sending_candidateBytes, &sending_candidate)

	if sending_candidateID == receiving_candidateID {
		// candidate voting to him/herself
		sending_candidate.Votes = sending_candidate.Votes + 1

	} else {
		receiving_candidateBytes, err := ctx.GetStub().GetState(receiving_candidateID)
		if err != nil {
			return fmt.Errorf("failed to read candidate state: %v", err)
		}
		if receiving_candidateBytes == nil {
			return fmt.Errorf("candidate %s does not exist", receiving_candidateID)
		}
		var receiving_candidate Candidate
		json.Unmarshal(receiving_candidateBytes, &receiving_candidate)

		if sending_candidate.Constituency != receiving_candidate.Constituency {
			return fmt.Errorf("voter:%s is not allowed to vote for candidate of another constituency", receiving_candidateID)
		}

		receiving_candidate.Votes = receiving_candidate.Votes + 1
		receiving_candidateBytes, _ = json.Marshal(receiving_candidate)
		err = ctx.GetStub().PutState(receiving_candidateID, receiving_candidateBytes)
		if err != nil {
			return fmt.Errorf("failed to update candidate votes: %v", err)
		}
	}

	sending_candidate.HasVoted = true
	sending_candidateBytes, _ = json.Marshal(sending_candidate)
	err = ctx.GetStub().PutState(sending_candidateID, sending_candidateBytes)
	if err != nil {
		return fmt.Errorf("failed to update voter status: %v", err)
	}

	return nil
}

/*
Get Result of all candidates

Get Result of by candidateID

Get Result of all constituency
   external database maintains list of candidates(id) constituency wise
     - call for each record by using candidateID
*/

func (s *SmartContract) GetResults(ctx contractapi.TransactionContextInterface) (map[string]int, error) {
	results := make(map[string]int)

	// Retrieve candidate count from the ledger (if stored)
	countBytes, err := ctx.GetStub().GetState("CANDIDATE_COUNT")
	if err != nil {
		return nil, fmt.Errorf("failed to read candidate count: %v", err)
	}

	if countBytes == nil {
		return nil, fmt.Errorf("failed: count isn't stored")
	}
	count, err := strconv.Atoi(string(countBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to convert count: %v", err)
	}

	// Iterate based on the stored count (assuming consecutive IDs)
	for i := 1; i <= count; i++ {
		candidateID := fmt.Sprintf("CAND%04d", i)
		candidateBytes, err := ctx.GetStub().GetState(candidateID)
		if err != nil {
			return nil, fmt.Errorf("failed to read candidate state for %s: %v", candidateID, err)
		}
		if candidateBytes == nil {
			continue // Skip if candidate doesn't exist (handles gaps)
		}

		var candidate Candidate
		json.Unmarshal(candidateBytes, &candidate)
		results[candidate.ID] = candidate.Votes

	}

	return results, nil
}

func (s *SmartContract) GetResultByID(ctx contractapi.TransactionContextInterface, candidateID string) (map[string]int, error) {
	results := make(map[string]int)

	candidateBytes, err := ctx.GetStub().GetState(candidateID)
	if err != nil {
		return nil, fmt.Errorf("failed to read candidate state: %v", err)
	}
	if candidateBytes == nil {
		return nil, fmt.Errorf("candidate %s does not exist", candidateID)
	}

	var candidate Candidate
	json.Unmarshal(candidateBytes, &candidate)
	results[candidate.ID] = candidate.Votes

	return results, nil
}

// Get Results Constituency wise
func (s *SmartContract) GetResultByVoterID(ctx contractapi.TransactionContextInterface, voterID string) ([]*Candidate, error) {
	// Step 1: Get voter from ledger
	voterBytes, err := ctx.GetStub().GetState(voterID)
	if err != nil {
		return nil, fmt.Errorf("failed to read voter state: %v", err)
	}
	if voterBytes == nil {
		return nil, fmt.Errorf("voter not found with ID %s", voterID)
	}

	var voter Voter
	err = json.Unmarshal(voterBytes, &voter)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal voter: %v", err)
	}

	// Step 2: Get candidate count
	countBytes, err := ctx.GetStub().GetState("CANDIDATE_COUNT")
	if err != nil {
		return nil, fmt.Errorf("failed to read candidate count: %v", err)
	}
	if countBytes == nil {
		return nil, fmt.Errorf("candidate count not found")
	}
	count, err := strconv.Atoi(string(countBytes))
	if err != nil {
		return nil, fmt.Errorf("invalid candidate count: %v", err)
	}

	// Step 3: Loop and collect candidates with matching constituency
	var candidates []*Candidate
	for i := 1; i <= count; i++ {
		candidateID := fmt.Sprintf("CAND%04d", i)
		candidateBytes, err := ctx.GetStub().GetState(candidateID)
		if err != nil {
			return nil, fmt.Errorf("failed to read candidate %s: %v", candidateID, err)
		}
		if candidateBytes == nil {
			continue
		}

		var candidate Candidate
		err = json.Unmarshal(candidateBytes, &candidate)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal candidate %s: %v", candidateID, err)
		}

		if candidate.Constituency == voter.Constituency {
			candidates = append(candidates, &candidate)
		}
	}

	return candidates, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create e-voting chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting e-voting chaincode: %v", err)
		return
	}
}
