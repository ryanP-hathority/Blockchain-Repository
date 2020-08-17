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

func initializeBallot() *Ballot {
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


func main() {

	ballot := initializeBallot()

	jsonstr, _ := json.Marshal(ballot)
	fmt.Println(string(jsonstr))

}

//REFRENCE: https://www.golangprograms.com/go-language/struct.html
