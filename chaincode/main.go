package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type TestChaincode struct {
}

// called whn cc is initialized; prepares ledger to all the processing
func (t *TestChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("~~~~~~~~ initializing chaincode ~~~~~~~~")
	return shim.Success(nil)
}

// all invoke requests are processed by this
func (t *TestChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("~~~~~~~~ ivoking chaincode ~~~~~~~~")

	function, args := stub.GetFunctionAndParameters()

	if function == "query" {
		return t.query(stub, args)
	} else if function == "transfer" {
		return t.invoke(stub, args)
	} else if function == "addAccount" {
		return t.addAccount(stub, args)
	}

	return shim.Error("Unknown action, check the first argument")
}

func (t *TestChaincode) addAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("~~~~~~~~~~~~~Adding account chaincode~~~~~~~~~~~~~")

	A := args[0]
	Aval, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("Added account successfully"))
}

// query
// Every readonly functions in the ledger will be here
func (t *TestChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("~~~~~~~~ querying ~~~~~~~~")

        AvalBal, err := stub.GetState(args[0])
        if err != nil {
            return shim.Error("can't get person1 state")
        }
        fmt.Println("~~~~~~~~~~~~~~~~~ BALANCE CHECK : PERSON 1 HAS %v",string(AvalBal))

        return shim.Success(AvalBal)
}

// invoke
// Every write function in the ledger will be here
func (t *TestChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("~~~~~~~~ invoking ~~~~~~~~")
        state1, err := stub.GetState(args[0])
        if err != nil {
            return shim.Error("can't get person1 state")
        }
        fmt.Println("~~~~~~~~~~~~~~~~~PERSON 1 HAS %v",string(state1))

        state2, err := stub.GetState(args[1])
        if err != nil {
            return shim.Error("can't get person2 state")
        }
        fmt.Println("~~~~~~~~~~~~~~~~~PERSON 2 HAS %v",string(state2))

        amt1, err := strconv.Atoi(string(state1))
        amt2, err := strconv.Atoi(string(state2))

        amount, err := strconv.Atoi(string(args[2]))
        fmt.Println("~~~~~~~~~~~~~~~~~TRYING TO TRANSFER %v",amount)

        new_amt1 := amt1 - amount
        new_amt2 := amt2 + amount
        fmt.Println("~~~~~~~~~~~~~~~~~NEW AMOUNT PERSON1 %v",new_amt1)
        fmt.Println("~~~~~~~~~~~~~~~~~NEW AMOUNT PERSON2 %v",new_amt2)
        if new_amt1<0 || new_amt2<0 {
            return shim.Error("not enough money, cannot overdraft")
        }
        err = stub.PutState(args[0], []byte(strconv.Itoa(int(new_amt1))))
	    if err != nil {
		    return shim.Error(err.Error())
	    }
        err = stub.PutState(args[1], []byte(strconv.Itoa(int(new_amt2))))
	    if err != nil {
		    return shim.Error(err.Error())
	    }
        fmt.Println("~~~~~~~~~~~~~~~~~TRANSFER SUCCESS")
        return shim.Success([]byte("transaction success"))
}


func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(TestChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}
