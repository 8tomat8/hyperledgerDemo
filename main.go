package main

import (
	"fmt"

	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Chaincode struct {
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// Init something...
	fmt.Println("Woohoo!!!")
	return shim.Success(nil)
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "createUser":
		return t.createUser(stub, args)
	case "transfer":
		return t.transfer(stub, args)
	case "getBalance":
		return t.getBalance(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"transfer\" \"getBalance\"")
}

func (t *Chaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	username := args[0]
	balance := args[1]

	err := stub.PutState(username, []byte(balance))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(username + " was created!"))
}

func (t *Chaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	from := args[0]
	to := args[1]
	amount, _ := strconv.Atoi(args[2])

	fromState, err := stub.GetState(from)
	if err != nil {
		return shim.Error(err.Error())
	}
	toState, err := stub.GetState(to)
	if err != nil {
		return shim.Error(err.Error())
	}

	fromStateInt, err := strconv.Atoi(string(fromState))
	if err != nil {
		return shim.Error(err.Error())
	}
	toStateInt, err := strconv.Atoi(string(toState))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(from, []byte(strconv.Itoa(fromStateInt-amount)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(to, []byte(strconv.Itoa(toStateInt+amount)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *Chaincode) getBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	username := args[0]
	data, err := stub.GetState(username)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(data)
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// CORE_PEER_ADDRESS=127.0.0.1:7051 CORE_CHAINCODE_ID_NAME=demo:0 go run main.go
// peer chaincode install -p github.com/8tomat8/hyperledgerDemo -n demo -v 0
// peer chaincode instantiate -n demo -v 0 -c '{"Args":["init"]}' -C mychannel

// peer chaincode invoke -n demo -v 0 -c '{"Args":["createUser", "user1", "100"]}' -C mychannel
// peer chaincode invoke -n demo -v 0 -c '{"Args":["createUser", "user2", "10"]}' -C mychannel
// peer chaincode invoke -n accountant -v 0 -c '{"Args":["transfer","user1","user2", "42"]}' -C mychannel
