package main

import (
	"fmt"
	"encoding/json"
)

type Candidate struct {
	candidateName string `json:"CandidateName"`
	isElected     bool   `json:"IsElected"`
}

type GovernmentRole struct {
	roleName         string      `json:"RoleName"`
	numEligable      int         `json:"NumEligable"`
	listOfCandidates []Candidate `json:"ListOfCandidates"`
}

type Ballot struct {
	ballotName  string           `json:"BallotName"`
	numOfRoles  int              `json:"NumOfRoles"`
	listOfRoles []GovernmentRole `json:"ListOfRoles"`
}


func main() {

	ballot := &Ballot{
		ballotName: "blockchain vote test",
		numOfRoles: 1,

		listOfRoles: []GovernmentRole{
			GovernmentRole{
				roleName: "U.S. President",
				numEligable: 1,

				listOfCandidates: []Candidate{
					Candidate {
						candidateName: "Ryan Patterson",
						isElected: true,
					},
				},
			},
		},
	}
	fmt.Println(ballot.ballotName)
	fmt.Println()
	fmt.Println(ballot.listOfRoles[0].numEligable)
	fmt.Println(ballot.listOfRoles[0].listOfCandidates[0].candidateName)

	jsonstr, _ := json.Marshal(ballot)
	fmt.Println(string(jsonstr))

}

//REFRENCE: https://www.golangprograms.com/go-language/struct.html
