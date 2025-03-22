package main

import (
	/*
		fmt: Format package
		encoding/json: to work with JSON
		hyperledger: contractapi
	*/

	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*
 SmartContract: It provides functions for managing votes. I have written in Go Language
 Voter struct represents a voter
 Candidate struct represents a candidate

 InitLedger initializes the ledger with some sample candidates (for testing)
 RegisterVoter registers a new voter
 CastVote : It allows a voter to cast their vote
 GetResults : It returns the vote count for each candidate
*/

type SmartContract struct {
	contractapi.Contract
}

type Voter struct {
	ID       string `json:"ID"`
	HasVoted bool   `json:"HasVoted"`
}

type Candidate struct {
	ID    string `json:"ID"`
	Name  string `json:"Name"`
	Votes int    `json:"Votes"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	candidates := []Candidate{
		{ID: "CAND001", Name: "Sahal", Votes: 0},
		{ID: "CAND002", Name: "Asim", Votes: 0},
		{ID: "CAND003", Name: "Jay", Votes: 0},
	}

	for _, candidate := range candidates {
		candidateBytes, _ := json.Marshal(candidate)
		err := ctx.GetStub().PutState(candidate.ID, candidateBytes)
		if err != nil {
			return fmt.Errorf("failed to put state for %s: %v", candidate.ID, err)
		}
	}
	return nil
}

func (s *SmartContract) RegisterVoter(ctx contractapi.TransactionContextInterface, voterID string) error {
	voter := Voter{
		ID:       voterID,
		HasVoted: false,
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

	candidate.Votes++
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

func (s *SmartContract) GetResults(ctx contractapi.TransactionContextInterface) (map[string]int, error) {
	results := make(map[string]int)

	// In a real application, you would iterate over all candidate IDs.
	// From database like MongoDB Cloud (Atlas) or Some Source
	// This example retrieves two hardcoded candidates.
	candidateIDs := []string{"CAND001", "CAND002"} // Replace with dynamic retrieval

	for _, candidateID := range candidateIDs {
		candidateBytes, err := ctx.GetStub().GetState(candidateID)
		if err != nil {
			return nil, fmt.Errorf("failed to read candidate state: %v", err)
		}
		if candidateBytes == nil {
			return nil, fmt.Errorf("candidate %s does not exist", candidateID)
		}

		var candidate Candidate
		json.Unmarshal(candidateBytes, &candidate)
		results[candidate.Name] = candidate.Votes
	}

	return results, nil
}

func main() {

	// CAD1: Sahal Pathan CAD2: Asim Pathan CAD3: Jay Pandya
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
