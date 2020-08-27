package chaincode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Vote struct {
	ID        string `json:"ID"`
	VoterName string `json:"voter"`
	Candidate string `json:"candidate"`
}

type Candidate struct {
        CandidateName string `json:"candidatename"`
        IsElected     bool   `json:"iselected"`
}

type GovernmentRole struct {
        RoleName         string      `json:"rolename"`
        NumEligable      int         `json:"numeligable"`
        ListOfCandidates []Candidate `json:"listofcandidates"`
}

type Ballot struct {
        BallotName  string           `json:"ballotname"`
        NumOfRoles  int              `json:"numofroles"`
	ID	    string           `json:"idfield"`
        ListOfRoles []GovernmentRole `json:"listofroles"`
}

func initializeBallot() *Ballot {

	ballot := &Ballot{
                BallotName: "blockchain vote test",
                NumOfRoles: 1,
                ID: "vote0",

                ListOfRoles: []GovernmentRole{
                        GovernmentRole{
                                RoleName: "U.S. President",
                                NumEligable: 1,

                                ListOfCandidates: []Candidate{
                                        Candidate {
                                                CandidateName: "Ryan Patterson",
                                                IsElected: false,
                                        },
					Candidate {
						CandidateName: "Andrew Bradjan",
						IsElected: false,
					},
					Candidate {
						CandidateName: "Saiteja Mallineni",
						IsElected: false,
					},
                                },
                        },
			GovernmentRole {
				RoleName: "U.S. Senate",
				NumEligable: 1,

				ListOfCandidates: []Candidate {
					Candidate {
						CandidateName: "Philip Bernick",
						IsElected: false,
					},
					Candidate {
						CandidateName: "Vishwam Annam",
						IsElected: false,
					},
				},
			},
			GovernmentRole {
				RoleName: "U.S. House of Representatives",
				NumEligable: 2,

				ListOfCandidates: []Candidate {
					Candidate {
						CandidateName: "Kiran Nagula",
						IsElected: false,
					},
					Candidate {
						CandidateName: "Rajeswari Gudiputi",
						IsElected: false,
					},
					Candidate {
						CandidateName: "Srija Vadem",
						IsElected: false,
					},
				},


			},
			GovernmentRole {
				RoleName: "State Senate",
				NumEligable: 1,

				ListOfCandidates: []Candidate {
					Candidate {
						CandidateName: "Tracie Gasik",
						IsElected: false,
					},
					Candidate {
						CandidateName: "Rhonda Steele",
						IsElected: false,
					},

				},
			},
			GovernmentRole {
				RoleName: "Court Attourney",
				NumEligable: 1,

				ListOfCandidates: []Candidate {
					Candidate {
						CandidateName: "Raj Katam",
						IsElected: false,
					},
					Candidate {
						CandidateName: "Sruthima Mantena",
						IsElected: false,
					},
				},
			},
			GovernmentRole {
				RoleName: "Court Sheriff",
				NumEligable: 1,

				ListOfCandidates: []Candidate {
					Candidate {
						CandidateName: "Aruna Ganta",
						IsElected: false,
					},
					Candidate {
						CandidateName: "Kristen Pond",
						IsElected: false,
					},
				},
			},
			GovernmentRole {
				RoleName: "County School Superindendent",
				NumEligable: 1,

				ListOfCandidates: []Candidate {
					Candidate {
						CandidateName: "Nagamani Kummari",
						IsElected: false,
					},
					Candidate {
						CandidateName: "Sena Reddy",
						IsElected: false,
					},
				},
			},
                },
        }

	return ballot
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	/*
	votes := []Vote{
		{ID: "vote1", VoterName: "andrew", Candidate: "bernie"},
		{ID: "vote2", VoterName: "ryan", Candidate: "bernie"},
		{ID: "vote3", VoterName: "saiteja", Candidate: "biden"},
		{ID: "vote4", VoterName: "philip", Candidate: "bernie"},
		{ID: "vote5", VoterName: "vishwam", Candidate: "biden"},
		{ID: "vote6", VoterName: "rhonda", Candidate: "warren"},
	}

	for _, vote := range votes {
		voteJSON, err := json.Marshal(vote)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(vote.ID, voteJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	*/
	//------   AUGMENTED CODE  --------
	ballot := initializeBallot()
	err = ctx.GetStub().PutState(ballot.ID, ballot)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateVote(ctx contractapi.TransactionContextInterface, string ballotAsJSON) error {

	/*
	exists, err := s.VoteExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the vote %s already exists", id)
	}

	vote := Vote{
		ID:        id,
		VoterName: voter,
		Candidate: candidate,
	}
	voteJSON, err := json.Marshal(vote)
	if err != nil {
		return err
	}

	*/
        //--------   AUGMENTED CODE----------
	var ballot Ballot
	err := Unmarshal(ballotAsJSON, &ballot)
	if(err != nil) {
		return err
	}
	exists, err := s.VoteExists(ctx, ballot.ID)
	if exists {
		return fmt.Errorf("The vote %s already exists", ballot.ID)
	}

	voteJSON, err := json.Marshal(ballot)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, voteJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadVote(ctx contractapi.TransactionContextInterface, id string) (*Vote, error) {
	voteJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if voteJSON == nil {
		return nil, fmt.Errorf("the vote %s does not exist", id)
	}

	var vote Vote
	err = json.Unmarshal(voteJSON, &vote)
	if err != nil {
		return nil, err
	}

	return &vote, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateVote(ctx contractapi.TransactionContextInterface, id string, voter string, candidate string) error {
	exists, err := s.VoteExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the vote %s does not exist", id)
	}

	// overwriting original asset with new asset
	vote := Vote{
		ID:        id,
		VoterName: voter,
		Candidate: candidate,
	}
	voteJSON, err := json.Marshal(vote)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, voteJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteVote(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.VoteExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the vote %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) VoteExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	voteJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return voteJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferVote(ctx contractapi.TransactionContextInterface, id string, newCandidate string) error {
	vote, err := s.ReadVote(ctx, id)
	if err != nil {
		return err
	}

	vote.Candidate = newCandidate
	voteJSON, err := json.Marshal(vote)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, voteJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllVotes(ctx contractapi.TransactionContextInterface) ([]*Vote, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var votes []*Vote
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var vote Vote
		err = json.Unmarshal(queryResponse.Value, &vote)
		if err != nil {
			return nil, err
		}
		votes = append(votes, &vote)
	}

	return votes, nil
}

func (s *SmartContract) TallyAllVotes(ctx contractapi.TransactionContextInterface) (string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Failed to get query results: " + err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	/*
		buffer, err := constructQueryResponseFromIteratorTallyAll(resultsIterator)
		if err != nil {
			return "", fmt.Errorf("Failed to construct query response: " + err.Error())
		}
	*/
	resultString := buffer.String()
	noBackSlashes := strings.Replace(resultString, "\\", "", -1)
	noleftBrackets := strings.Replace(noBackSlashes, `{`, ``, -1)
	norightBrackets := strings.Replace(noleftBrackets, `}`, ``, -1)
	noQuotes := strings.Replace(norightBrackets, `"`, ``, -1)
	splitString := strings.Split(noQuotes, ",")
	fmt.Println(" - Begin iterating through votes")
	var candidates []string
	i := 0
	for range splitString {
		if strings.HasPrefix(splitString[i], `candidate:`) == true {
			found := contains(candidates, splitString[i])
			if found == false {
				candidates = append(candidates, splitString[i])
				i++
			} else {
				i++
			}
		} else {
			i++
		}
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("No votes in blockchain to tally")
	}
	fmt.Println(" - Begin tallying votes for each candidate")
	i = 0
	var candidatesTally []string
	newString := fmt.Sprint(splitString)
	for range candidates {
		count := strings.Count(newString, candidates[i])
		candidatesTally = append(candidatesTally, candidates[i]+" - vote total: "+strconv.Itoa(count))
		i++
	}
	strCandidatesTally := strings.Join(candidatesTally, "\n")
	result := strings.Join(candidatesTally, " ")
	fmt.Println(" - End tallying votes")
	fmt.Println(strCandidatesTally)
	return result, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

