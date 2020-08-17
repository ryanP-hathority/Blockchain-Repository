package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	//"strings"
	//"time"
	"os"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// vote implements a chaincode to manage a vote
type SmartContract struct {
	contractapi.Contract
}

type Vote struct {
	VoterName string `json:"voter"`
	Candidate string `json:"candidate"`
}

type voteList struct {
	KeyValue string
	Vote     Vote
}

//Init is called during chaincode instantiation to initialize any
//data. Note that chaincode upgrade also calls this function to reset
//or to migrate data.
// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	votes := []Vote{
		Vote{VoterName: "Andrew", Candidate: "Joe"},
		Vote{VoterName: "Ryan", Candidate: "Ben"},
		Vote{VoterName: "Saiteja", Candidate: "Kamala"},
	}

	for i, vote := range votes {
		voteAsBytes, _ := json.Marshal(vote)
		err := ctx.GetStub().PutState("VOTE"+strconv.Itoa(i), voteAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}
/*
// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
//peer chaincode invoke -n mycc -c '{"Args":["set", "voterID", "voterName", "candidate"]}' -C myc
func (v *SmartContract) Invoke(ctx contractapi.TransactionContextInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := ctx.GetStub().GetFunctionAndParameters()

	var result string
	var err error
	if fn == "insert" {
		result, err = insert(ctx, args)
	} else if fn == "delete" {
		result, err = remove(ctx, args)
	} else if fn == "tallyAll" {
		result, err = tallyAll(ctx, args)
	} else if fn == "tallyForcandidate" {
		result, err = tallyForcandidate(ctx, args)
	} else if fn == "changeVote" {
		result, err = changeVote(ctx, args)
	} else if fn == "getVoterscandidate" {
		result, err = getVoterscandidate(ctx, args)
	} else if fn == "queryByID" {
		result, err = queryByID(ctx, args)
	} else if fn == "getHistoryForVote" {
		result, err = getHistoryForVote(ctx, args)
	} else {
		fmt.Println("Invalid command")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}
*/
//insert creates a new vote and stores it into the chaincode state.
func (v *SmartContract) insert(ctx contractapi.TransactionContextInterface, voterKey string, votername string, candidate string,) error {
	vote := Vote{
		VoterName:   votername,
		Candidate:  candidate,
	}

	voteAsBytes, _ := json.Marshal(vote)

	return ctx.GetStub().PutState(voterKey, voteAsBytes)

}

func verifyIdentity() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("vote")

	result, err := contract.EvaluateTransaction("queryAllCars")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = contract.SubmitTransaction("createCar", "CAR10", "VW", "Polo", "Grey", "Mary")
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	_, err = contract.SubmitTransaction("changeCarOwner", "CAR10", "Archie")
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}

	result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create vote chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting vote chaincode: %s", err.Error())
	}
}

