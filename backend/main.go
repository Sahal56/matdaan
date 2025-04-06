package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"evoting_server/config"
	"evoting_server/handlers"
	"evoting_server/services"

	"github.com/gin-gonic/gin"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "/home/ubuntu/Hyperledger/matdaan/hyperledger-fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

func main() {
	// gRPC client connection for Hyperledger Fabric
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	chaincodeName := "evoting"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB()

	awsClients, err := config.InitAWS()
	if err != nil {
		log.Fatal(err)
	}

	// votingService := services.NewVotingService(db.Voters, db.Candidates, contract)
	votingService := services.NewVotingService(
		db.Voters,
		db.Candidates,
		contract,
		awsClients.S3,
		awsClients.Rekognition,
		awsClients.BucketName,
	)
	votingHandler := handlers.NewHandler(votingService)

	router := gin.Default()

	// CORS
	// router.Use(cors.Default()) // Enable CORS with default settings | dev testing

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Routes
	api := router.Group("/api")
	{
		api.POST("/candidates", votingHandler.GetCandidatesByVoter)
		api.POST("/register", votingHandler.RegisterVoter)

		api.POST("/login", votingHandler.LoginVoter)
		api.GET("/me", votingHandler.GetCurrentVoter)

		api.POST("/vote", votingHandler.CastVote)
		api.POST("/resultsByVoter", votingHandler.GetResultsByVoter)
		api.GET("/results", votingHandler.GetResults)

	}

	// Protected routes which requires Login
	// protected := router.Group("/api")
	// protected.Use(AuthMiddleware())
	// {
	// 	protected.POST("/vote", votingHandler.CastVote)
	// 	protected.POST("/resultsByVoter", votingHandler.GetResultsByVoter)
	// }

	// Testing
	//  Python Server python3 -m http.server

	// Server Ports
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on :" + port)
	router.Run(":" + port)
}

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		voterID, err := c.Cookie("voterID")
// 		if err != nil || voterID == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
// 			c.Abort()
// 			return
// 		}
// 		c.Next()
// 	}
// }

func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}
