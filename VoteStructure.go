package main

import (
	"fmt"
	"encoding/json"
)

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
	ListOfRoles []GovernmentRole `json:"listofroles"`
}


func main() {

	ballot := &Ballot{
		BallotName: "blockchain vote test",
		NumOfRoles: 1,

		ListOfRoles: []GovernmentRole{
			GovernmentRole{
				RoleName: "U.S. President",
				NumEligable: 1,

				ListOfCandidates: []Candidate{
					Candidate {
						CandidateName: "Ryan Patterson",
						IsElected: true,
					},
				},
			},
		},
	}
	fmt.Println(ballot.BallotName)
	fmt.Println()
	fmt.Println(ballot.ListOfRoles[0].NumEligable)
	fmt.Println(ballot.ListOfRoles[0].ListOfCandidates[0].CandidateName)

	jsonstr, _ := json.Marshal(ballot)
	fmt.Println(string(jsonstr))

}

//REFRENCE: https://www.golangprograms.com/go-language/struct.html
