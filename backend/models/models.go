package models

type Voter struct {
	ID           string `bson:"_id,omitempty" json:"ID"`
	Name         string `bson:"name" json:"Name" binding:"required"`
	Constituency string `bson:"constituency" json:"Constituency"`
	Age          int    `bson:"age" json:"Age"`
	HasVoted     bool   `bson:"hasVoted" json:"HasVoted"`
	// PhotoPath    string `bson:"photopath" json:"PhotoPath" binding:"required"`
}

type Candidate struct {
	ID           string `bson:"_id,omitempty" json:"ID"`
	Name         string `bson:"name" json:"Name" binding:"required"`
	Constituency string `bson:"constituency" json:"Constituency" binding:"required"`
	Votes        int    `bson:"votes" json:"Votes" binding:"required"`
	HasVoted     bool   `bson:"hasVoted" json:"HasVoted" binding:"required"`
}

// VoteRequest : request body for casting a vote
type VoteRequest struct {
	VoterID     string `json:"voter_id"`
	CandidateID string `json:"candidate_id"`
}

// CandidateVoteRequest : request body for a candidate voting
type CandidateVoteRequest struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
}
