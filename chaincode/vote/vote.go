package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"strings"
	"strconv"
	"time"

        "github.com/hyperledger/fabric/core/chaincode/shim"
        "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// vote implements a chaincode to manage a vote
type SmartContract struct {
	contractapi.Contract
}

type Vote struct {
	VoterName  string `json:"voter"`
	Candidate  string `json:"candidate"`
}

type voteList struct {
	KeyValue string
	Vote Vote
}

//Init is called during chaincode instantiation to initialize any
//data. Note that chaincode upgrade also calls this function to reset
//or to migrate data.
func (v *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println(" - Inititializing the chaincode")
	// Get the args from the transaction proposal1
	args := stub.GetStringArgs()
        if len(args) != 2 {
                return shim.Error("Incorrect arguments. Expecting a voter's name and a candidate")
        }

        // Set up any variables or assets here by calling stub.PutState()

	//formed args into bytes using the json function

	voterName        := strings.ToLower(args[0])
        candidateName    := strings.ToLower(args[1])
        //var attributesForKey []string
        //attributesForKey = append(attributesForKey, voterName)
        //attributesForKey = append(attributesForKey, candidateName)
        //key, err         := stub.CreateCompositeKey("vote", attributesForKey)
	
	//USING TXID:::::
	key := stub.GetTxID()
        //if err != nil {
        //        return shim.Error("Failed to create key: " + err.Error())
        //}

	vote := Vote{voterName, candidateName}
	voteJSONAsBytes, err := json.Marshal(vote)

        if err != nil {
                return shim.Error("Failed to construct vote: " + err.Error())
        }

        // We store the key and the value on the ledger
        err = stub.PutState(key, voteJSONAsBytes)
        if err != nil {
                return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
        }

	fmt.Printf(" - end initializing chaincode\n")
        fmt.Println("   key: " + key)
        fmt.Println("   voter: " + voterName)
        fmt.Println("   candidate: " + candidateName)

        return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
//peer chaincode invoke -n mycc -c '{"Args":["set", "voterID", "voterName", "candidate"]}' -C myc
func (v *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
        fn, args := stub.GetFunctionAndParameters()

        var result string
        var err error
        if fn == "insert" {
                result, err = insert(stub, args)
        } else if fn == "delete" {
                result, err = remove(stub, args)
        } else if fn == "tallyAll" {
		result, err = tallyAll(stub, args)
	} else if fn == "tallyForcandidate" {
		result, err = tallyForcandidate(stub, args)
	} else if fn == "changeVote" {
		result, err = changeVote(stub, args)
	} else if fn == "getVoterscandidate" {
		result, err = getVoterscandidate(stub, args)
	} else if fn == "queryByID" {
		result, err = queryByID(stub, args)
	} else if fn == "getHistoryForVote" {
		result, err = getHistoryForVote(stub, args)
	} else {
		fmt.Println("Invalid command")
	}

        if err != nil {
                return shim.Error(err.Error())
        }

        // Return the result as success payload
        return shim.Success([]byte(result))
}

//insert creates a new vote and stores it into the chaincode state.
func insert(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var result string

	if len(args) != 2 {
		return "", fmt.Errorf("Expecting two arguments. Both voter and candidate are required")
	}
	fmt.Println(" - Begin initializing vote")

	// verify inputs
	if len(args[0]) <= 0 {
		return "", fmt.Errorf("Expecting a non-empty string for the voter field")
	}
	if len(args[1]) <= 0 {
		return "", fmt.Errorf("Expecting a non-empty string for the candidate field")
	}

	voterName        := strings.ToLower(args[0])
	candidateName    := strings.ToLower(args[1])
	var attributesForKey []string
	attributesForKey = append(attributesForKey, voterName)
        attributesForKey = append(attributesForKey, candidateName)
	//OLD KEY 
	key, err         := stub.CreateCompositeKey("vote", attributesForKey)


	//USING TXID:::::
	key = stub.GetTxID()
	if err != nil {
		return "", fmt.Errorf("Failed to create key: " + err.Error())
	}
	//TODO: possibly need to verify the voter and candidate do not exist in the ledger

	vote := Vote{VoterName: voterName, Candidate: candidateName}
	voteJSONAsBytes, err := json.Marshal(vote)
	if err != nil {
		return "", fmt.Errorf("Failed to construct vote: " + err.Error())
	}

	err = stub.PutState(key, voteJSONAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to insert vote: " + err.Error())
	}

	fmt.Println(" - end initializing vote")
	fmt.Println("	key: " + key)
	fmt.Println("	voter: " + voterName)
	fmt.Println("	candidate: " + candidateName)
	result = "Created Vote for voter " + voterName + " with key: " + key

	return result, nil
}

//remove removes a vote key/value pair from the chaincode state.
func remove(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var result string
        if len(args) != 1 {
                return "", fmt.Errorf("Incorrect number of arguments. Expecting a voter ID")
        }
        fmt.Println(" - Begin deleting vote")

        if len(args[0]) <= 0 {
                return "", fmt.Errorf("Expecting a non-empty string for the voter ID field")
        }

        voterID := args[0]
	keyvalue, err := stub.GetState(voterID)
	if keyvalue == nil {
		return "", fmt.Errorf("VoterID was not found in the blockchain")
	}
        err = stub.DelState(voterID)
        if err != nil {
                return "", fmt.Errorf("Failed to delete vote: ", voterID, err)
        }

        fmt.Println(" - End deleting vote")
        result = "Deleted Vote"

        return result, nil
}

// tallyAll parses the chaincode and returns a string value of how many votes 
// there are for each candidate. 
func tallyAll(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	resultsIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("Failed to get query results: " + err.Error())
	}
	defer resultsIterator.Close()

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
			found:= contains(candidates, splitString[i])
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
		candidatesTally = append(candidatesTally, candidates[i] + " - vote total: " + strconv.Itoa(count))
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
func tallyForcandidate(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting a candidate.")
	}
	
	resultsIterator, err := stub.GetStateByRange("", "")
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
	count := strings.Count(resultingString, `"candidate":"` + candidate)
	return "- tallyForcandidate results: " + string(count), nil
}


//changeVote changes a vote by setting a new candidate name on the vote.
func changeVote(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	//var key []byte
        var result2 string

        //      0          1                   2
        // "voterID", "Andrew Bradjan", "Bernie Sanders"
        if len(args) != 3 {
                return "", fmt.Errorf("Incorrect number of arguments. Expecting a voter ID and a new Candidate")
        }
        fmt.Println(" - Begin changing vote")

        if len(args[0]) <= 0 {
                return "", fmt.Errorf("Expecting a non-empty string for the voter ID field")
        }
        if len(args[1]) <= 0 {
                return "", fmt.Errorf("Expecting a non-empty string for the voter name field")
        }
        if len(args[2]) <= 0 {
              return "", fmt.Errorf("Expecting a non-empty string for the new candidate field")
        }
        voterID := args[0]
        voterName := strings.ToLower(args[1])
        newCandidate := strings.ToLower(args[2])

        key, err := stub.GetState(voterID)
        if key != nil {
                return "", fmt.Errorf("Failed to get vote:" + voterID)
        }
        if err != nil {
                return "", fmt.Errorf("Failed to get vote:" + voterID)
        }

        vote := Vote{VoterName: voterName, Candidate: newCandidate}
        voteJSONAsBytes, err := json.Marshal(vote)
        if err != nil {
                return "", fmt.Errorf("Failed to construct new vote: " + err.Error())
        }
        //err = stub.PutState(key, voteJSONAsBytes)
        //if err != nil {

        err = stub.PutState(voterID, voteJSONAsBytes) //rewrite the vote
        if err != nil {
              return "", fmt.Errorf("Failed to change vote:" + err.Error())
        }
        fmt.Println(" - end changeVote")
        result2 = "Changed Vote"
        return result2, nil
}

func getVoterscandidate(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var result string

	return result, nil
}

//queryByID reads a vote from the chaincode state by the ID.
func queryByID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting a voter ID")
	}
	
	//key := args[0]
	//fmt.Println("Querying chaincode for key: " + key)
	result, err := stub.GetState(args[0])
	
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if result == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}

	fmt.Println("result: " + string(result))

	return string(result), nil
}

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
func getHistoryForVote(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting a voter ID")
	}
	fmt.Printf("- start getHistoryForVote: %s\n", args[0])

	resultsIterator, err := stub.GetHistoryForKey(args[0])
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

/*
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
        if len(args) != 2 {
                return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
        }

        err := stub.PutState(args[0], []byte(args[1]))
        if err != nil {
                return "", fmt.Errorf("Failed to set asset: %s", args[0])
        }
        return args[1], nil
}

*/
// Get returns the value of the specified asset key
func getValue(stub shim.ChaincodeStubInterface, args []string) (string, error) {
        if len(args) != 1 {
                return "", fmt.Errorf("Incorrect arguments. Expecting a key")
        }

        value, err := stub.GetState(args[0])
        if err != nil {
                return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
        }
        if value == nil {
                return "", fmt.Errorf("Asset not found: %s", args[0])
        }
        return string(value), nil
}


// main function starts up the chaincode in the container during instantiate
func main() {
        if err := shim.Start(new(SimpleChaincode)); err != nil {
                fmt.Printf("Error starting Simple chaincode: %s", err)
        }
}



















