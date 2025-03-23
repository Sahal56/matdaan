package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*
operations:
    Register Voter
	Cast Vote


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
	ID           string `json:"ID"`
	Constituency string `json:"Constituency"`
	Votes        int    `json:"Votes"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// statically giving input
	candidates := []Candidate{
		{ID: "CAND001", Constituency: "Anand", Votes: 0},
		{ID: "CAND002", Constituency: "Anand", Votes: 0},
		{ID: "CAND003", Constituency: "Anand", Votes: 0},
		{ID: "CAND004", Constituency: "Vadodara", Votes: 0},
		{ID: "CAND005", Constituency: "Vadodara", Votes: 0},
		{ID: "CAND006", Constituency: "Vadodara", Votes: 0},
	}

	// Dynamically Candidates can be filled based on particular channel's peer nodes based on constituency

	for _, candidate := range candidates {
		candidateBytes, _ := json.Marshal(candidate)
		err := ctx.GetStub().PutState(candidate.ID, candidateBytes)
		if err != nil {
			return fmt.Errorf("failed to put state for %s: %v", candidate.ID, err)
		}
	}
	return nil
}

func (s *SmartContract) RegisterVoter(ctx contractapi.TransactionContextInterface,
	voterID string, constituency string) error {

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

func (s *SmartContract) GetResults(ctx contractapi.TransactionContextInterface) (map[string]int, error) {
	results := make(map[string]int)

	// In a real application, you would iterate over all candidate IDs.
	// From database like MongoDB Cloud (Atlas) or Some Source
	// This example retrieves two hardcoded candidates.

	// It can have dynamic values
	candidateIDs := []string{"CAND001", "CAND002"}

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
