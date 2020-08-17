package main

import (
	"fmt"
	//"strings"
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

type Ballot struct {
	ballotName string
	numOfRoles int
	listOfRoles []GovernmentRole
}


func main() {
	var candidate = Candidate {
		candidateName: "Ryan Patterson",
		isElected: true,
	}

	fmt.Println(candidate.candidateName)
	if (candidate.isElected) {
		fmt.Println("elected")
	}

	var role = GovernmentRole {
		roleName: "U.S. President",
		numEligable: 1,
		listOfCandidates: []Candidate {
			Candidate {
				candidateName: "Ryan Patterson",
				isElected: true,
			},
		},
	}

	fmt.Println(role.roleName)
	fmt.Println(role.numEligable)
	fmt.Println(role.listOfCandidates[0].candidateName)




	var ballot = Ballot{
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

}

//REFRENCE: https://www.golangprograms.com/go-language/struct.html
