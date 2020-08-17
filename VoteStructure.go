package main

import (
	"fmt"
	"strings"
)

type Candidate struct {
	candidateName string
	isElected bool
}

type GovernmentRole struct {
	roleName string
	numEligable int
	listOfCandidates []Candidate
}

type Ballot stuct {
	ballotName string
	numOfRoles int
	listOfRoles []GovernmentRole
}


func main() {
	ballot := Ballot{
		ballotName: "blockchain vote test",
		numOfRoles: 1,

		listOfRoles: []GovernmentRole{
			roleName: "U.S. President",
			numEligable: 1,

			listOfCandidates []Candidate{
				candidateName: "Ryan Patterson",
				isElected: true,
			},
		},
	}
	fmt.Println(ballot.ballotName)
	fmt.Println()
	fmt.Println(ballot.listOfRoles[0].numOfRoles)
}

//REFRENCE: https://www.golangprograms.com/go-language/struct.html
