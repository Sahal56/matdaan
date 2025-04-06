package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"evoting_server/services"

	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	Service *services.VotingService
}

// NewHandler initializes a new handler
func NewHandler(service *services.VotingService) *Handler {
	return &Handler{Service: service}
}

// LoginVoter handles voter login and face verification
func (h *Handler) LoginVoter(c *gin.Context) {
	var req struct {
		VoterID string `json:"voterID"`
		Photo   string `json:"photo"` // Base64 encoded image
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Invalid JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	fmt.Println("Parsed input:", req.VoterID)

	// Decode base64 photo
	base64Data := req.Photo
	if commaIdx := strings.Index(base64Data, ","); commaIdx != -1 {
		base64Data = base64Data[commaIdx+1:]
	}
	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		fmt.Println("Failed to decode image:", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid image data"})
		return
	}

	// Verify voter
	success, err := h.Service.VerifyVoter(req.VoterID, imageBytes)
	if err != nil || !success {
		fmt.Println("Verification failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed"})
		return
	}

	fmt.Println("Authentication success")
	c.SetCookie("voterID", req.VoterID, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) GetCurrentVoter(c *gin.Context) {
	voterID, err := c.Cookie("voterID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not logged in"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"voterID": voterID})
}

type RegisterVoterRequest struct {
	VoterID      string `json:"voterID" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Constituency string `json:"constituency" binding:"required"`
	Age          int    `json:"age" binding:"required"`
}

func (h *Handler) GetCandidatesByVoter(c *gin.Context) {
	var req struct {
		VoterID string `json:"voterID"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.VoterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "voterID is required"})
		return
	}

	if h.Service == nil {
		log.Println("Service is nil!")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	candidates, err := h.Service.GetCandidatesByVoterService(c.Request.Context(), req.VoterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"candidates": candidates})
}

// RegisterVoter handles voter registration with image upload
func (h *Handler) RegisterVoter(c *gin.Context) {
	voterID := c.PostForm("voterID")
	name := c.PostForm("name")
	constituency := c.PostForm("constituency")
	ageStr := c.PostForm("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid age"})
		return
	}

	// Read image file
	file, _, err := c.Request.FormFile("faceImage")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Image file is required"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to read image"})
		return
	}

	// Call service
	err = h.Service.RegisterVoter(
		c.Request.Context(),
		voterID,
		name,
		constituency,
		age,
		imageBytes,
	)

	if err != nil {
		if err.Error() == "voter already registered" || err.Error() == "voter must be at least 18 years old" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Server error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Voter registered successfully"})
}

// CastVote handles vote casting
func (h *Handler) CastVote(c *gin.Context) {
	var request struct {
		VoterID     string `json:"voterID"`
		CandidateID string `json:"candidateID"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.CastVote(request.VoterID, request.CandidateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote cast successfully"})
}

// GetResults returns election results
func (h *Handler) GetResults(c *gin.Context) {
	results, err := h.Service.GetResults()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *Handler) GetResultsByVoter(c *gin.Context) {
	var payload struct {
		VoterID string `json:"voterID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "voterID is required"})
		return
	}

	results, err := h.Service.GetResultsByVoterService(payload.VoterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
