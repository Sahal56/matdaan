package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AWSClients struct {
	S3          *s3.Client
	Rekognition *rekognition.Client
	BucketName  string
}

// Database holds the MongoDB client and collections
type Database struct {
	Client     *mongo.Client
	DB         *mongo.Database
	Voters     *mongo.Collection
	Candidates *mongo.Collection
}

func InitAWS() (*AWSClients, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		return nil, fmt.Errorf("AWS_S3_BUCKET environment variable not set")
	}

	return &AWSClients{
		S3:          s3.NewFromConfig(cfg),
		Rekognition: rekognition.NewFromConfig(cfg),
		BucketName:  bucket,
	}, nil
}

func InitDB() (*Database, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using existing environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI not set in environment")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database with context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")

	// Initialize collections
	db := client.Database("evotingDB")
	return &Database{
		Client:     client,
		DB:         db,
		Voters:     db.Collection("voters"),
		Candidates: db.Collection("candidates"),
	}, nil
}

// CloseDB disconnects from MongoDB
func (db *Database) CloseDB() {
	if db.Client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Client.Disconnect(ctx); err != nil {
		log.Println("Error closing MongoDB connection:", err)
	} else {
		fmt.Println("Disconnected from MongoDB")
	}
}
