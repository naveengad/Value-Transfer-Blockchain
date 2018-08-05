package main

import (
	"fmt"
	"balance_transfer/blockchain"
	"os"
)

func main() {
	fSetup := blockchain.FabricSetup{
	OrgAdmin:        "Admin",
	OrgName:         "Org1",
    OtherOrg:        "Org2",
    ConfigFile:      "config.yaml",
	ChannelID:       "channel1",
	ChannelConfig:   os.Getenv("GOPATH") + "/src/balance_transfer/network/artifacts/baltransfer.channel1.tx",
	ChainCodeID:     "cc1",
    ChaincodePath:   "balance_transfer/chaincode/",
    ChaincodeGoPath: os.Getenv("GOPATH"),
    UserName:        "User1",
    }

    // Initialization of the Fabric SDK from the previously set properties
    err := fSetup.Initialize()
    if err != nil {
	fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
    }

    err = fSetup.InstallAndInstantiateCC()
    if err != nil {
        fmt.Printf("install and instantiate failed %v\n", err)
    }

    var accArgs1 []string
    accArgs1 = append(accArgs1, "p1")
    accArgs1 = append(accArgs1, "100")
    response, err := fSetup.QueryAddAccountTest(accArgs1)
    if err != nil {
        fmt.Printf("TRANSFER FAIL error querying %v\n",err)
    } else {
        fmt.Printf("p1: %s\n",response)
    }

    var accArgs2 []string
    accArgs2 = append(accArgs2, "p2")
    accArgs2 = append(accArgs2, "50")
    response, err = fSetup.QueryAddAccountTest(accArgs2)
    if err != nil {
        fmt.Printf("TRANSFER FAIL error querying %v\n",err)
    } else {
        fmt.Printf("p2: %s\n",response)
    }

    response, err = fSetup.QueryVerifyTest("p1")
    if err != nil {
        fmt.Printf("error querying %v\n",err)
    } else {
        fmt.Printf("~~~~~~~~~~~~~~~~~~~~~BALANCE QUERY for p1 before transfer is %s\n",response)
    }

    response, err = fSetup.QueryVerifyTest("p2")
    if err != nil {
        fmt.Printf("error querying %v\n",err)
    } else {
        fmt.Printf("~~~~~~~~~~~~~~~~~~~~~BALANCE QUERY for p2 before transfer is %s\n",response)
    }

    var trnsfrArgs []string
    trnsfrArgs = append(trnsfrArgs, "p1")
    trnsfrArgs = append(trnsfrArgs, "p2")
    trnsfrArgs = append(trnsfrArgs, "50")
    response, err = fSetup.QueryTransactionTest(trnsfrArgs)
    if err != nil {
        fmt.Printf("TRANSFER FAIL error querying %v\n",err)
    } else {
        fmt.Printf("TRANSFER SUCCESS response from query is %s\n",response)
    }

    response, err = fSetup.QueryVerifyTest("p1")
    if err != nil {
        fmt.Printf("error querying %v\n",err)
    } else {
        fmt.Printf("~~~~~~~~~~~~~~~~~~~~~BALANCE QUERY for p1 after transfer is %s\n",response)
    }

    response, err = fSetup.QueryVerifyTest("p2")
    if err != nil {
        fmt.Printf("error querying %v\n",err)
    } else {
        fmt.Printf("~~~~~~~~~~~~~~~~~~~~~BALANCE QUERY for p2 after transfer is %s\n",response)
    }
}
