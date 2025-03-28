package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	// "github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
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
	ID           string `json:"ID"`
	Constituency string `json:"Constituency"`
	Votes        int    `json:"Votes"`
	HasVoted     bool   `json:"HasVoted"` // new logic: candidates can also vote
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// statically giving input
	candidates := []Candidate{
		{ID: "CAND0001", Constituency: "Anand", Votes: 0},
		{ID: "CAND0002", Constituency: "Anand", Votes: 0},
		{ID: "CAND0003", Constituency: "Anand", Votes: 0},
		{ID: "CAND0004", Constituency: "Vadodara", Votes: 0},
		{ID: "CAND0005", Constituency: "Vadodara", Votes: 0},
		{ID: "CAND0006", Constituency: "Vadodara", Votes: 0},
	}

	// Counting the number of candidates
	numCandidates := len(candidates)

	// Dynamically Candidates can be filled based on particular channel's peer nodes based on constituency

	for _, candidate := range candidates {
		candidateBytes, _ := json.Marshal(candidate)
		err := ctx.GetStub().PutState(candidate.ID, candidateBytes)
		if err != nil {
			return fmt.Errorf("failed to put state for %s: %v", candidate.ID, err)
		}
	}

	// storing numCandidates
	err := ctx.GetStub().PutState("CANDIDATE_COUNT", []byte(strconv.Itoa(numCandidates)))
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

	// In a real application, you would iterate over all candidate IDs.
	// From database like MongoDB Cloud (Atlas) or Some Source
	// This example retrieves two hardcoded candidates.
	// It can have dynamic values
	// candidateIDs := []string{"CAND001", "CAND002"} // as of know only 2 values

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
