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
/*
//remove removes a vote key/value pair from the chaincode state.
func (v *SmartContract) remove(ctx contractapi.TransactionContextInterface, args []string) error {
	var result string
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting a voter ID")
	}
	fmt.Println(" - Begin deleting vote")

	if len(args[0]) <= 0 {
		return "", fmt.Errorf("Expecting a non-empty string for the voter ID field")
	}

	voterID := args[0]
	keyvalue, err := ctx.GetStub().GetState(voterID)
	if keyvalue == nil {
		return "", fmt.Errorf("VoterID was not found in the blockchain")
	}
	err = ctx.GetStub().DelState(voterID)
	if err != nil {
		return "", fmt.Errorf("Failed to delete vote: ", voterID, err)
	}

	fmt.Println(" - End deleting vote")
	result = "Deleted Vote"

	return result, nil
}

// tallyAll parses the chaincode and returns a string value of how many votes
// there are for each candidate.
func (v *SmartContract) tallyAll(ctx contractapi.TransactionContextInterface, args []string) (string, error) {
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
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	buffer, err := constructQueryResponseFromIteratorTallyAll(resultsIterator)
	if err != nil {
		return "", fmt.Errorf("Failed to construct query response: " + err.Error())
	}

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
*/
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
/*
func constructQueryResponseFromIteratorTallyAll(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	return &buffer, nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// tallyForcandidate parses the chaincode and returns a string value of how many votes
// there are for a specific candidate.
func (v *SmartContract) tallyForcandidate(ctx contractapi.TransactionContextInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting a candidate.")
	}

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Failed to get query results: " + err.Error())
	}
	defer resultsIterator.Close()
	candidate := strings.ToLower(args[0])
	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return "", fmt.Errorf("Failed to construct query response: " + err.Error())
	}
	resultingString := buffer.String()
	count := strings.Count(resultingString, `"candidate":"`+candidate)
	return "- tallyForcandidate results: " + string(count), nil
}
*/
//changeVote changes a vote by setting a new candidate name on the vote.
func (v *SmartContract) ChangeCandidate(ctx contractapi.TransactionContextInterface, voterKey string, newCandidate string) error {
	vote, err := v.QueryVote(ctx, voterKey)

	if err != nil {
		return err
	}

	vote.Candidate = newCandidate

	voteAsBytes, _ := json.Marshal(vote)

	return ctx.GetStub().PutState(voterKey, voteAsBytes)
}
/*
func (v *SmartContract) getVoterscandidate(ctx contractapi.TransactionContextInterface, args []string) (string, error) {
	var result string

	return result, nil
}
*/
//queryVote reads a vote from the chaincode state by the ID.
func (v *SmartContract) QueryVote(ctx contractapi.TransactionContextInterface, voterKey string) (*Vote, error) {
	voteAsBytes, err := ctx.GetStub().GetState(voterKey)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if voteAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", voterKey)
	}

	vote := new(Vote)
	_ = json.Unmarshal(voteAsBytes, vote)

	return vote, nil
}
/*
func constructHistoryResponseFromIterator(resultsIterator shim.HistoryQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		if queryResponse.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(queryResponse.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(queryResponse.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// getHistoryForVote returns the history of the specified asset key
func (v *SmartContract) getHistoryForVote(ctx contractapi.TransactionContextInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting a voter ID")
	}
	fmt.Printf("- start getHistoryForVote: %s\n", args[0])

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get query results: " + err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := constructHistoryResponseFromIterator(resultsIterator)
	if err != nil {
		return "", fmt.Errorf("Failed to construct query response: " + err.Error())
	}

	return "- getHistoryForVote queryResult: " + buffer.String(), nil
}
*/
/*
func set(ctx contractapi.TransactionContextInterface, args []string) (string, error) {
        if len(args) != 2 {
                return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
        }
        err := ctx.GetStub().PutState(args[0], []byte(args[1]))
        if err != nil {
                return "", fmt.Errorf("Failed to set asset: %s", args[0])
        }
        return args[1], nil
}
*/
/*
// Get returns the value of the specified asset key
func (v *SmartContract) getValue(ctx contractapi.TransactionContextInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := ctx.GetStub().GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}
*/
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

