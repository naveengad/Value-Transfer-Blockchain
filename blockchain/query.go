package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
)

func (setup *FabricSetup) QueryTransactionTest(args []string) (string, error) {

	response, err := setup.client1.Execute(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: "transfer", Args: [][]byte{[]byte(args[0]),[]byte(args[1]),[]byte(args[2])}})

	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}

func (setup *FabricSetup) QueryVerifyTest(name string) (string, error) {

	response, err := setup.client2.Query(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: "query", Args: [][]byte{[]byte(name)}})

	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}

func (setup *FabricSetup) QueryAddAccountTest(args []string) (string, error) {

	response, err := setup.client1.Execute(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: "addAccount", Args: [][]byte{[]byte(args[0]),[]byte(args[1])}})

	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}
