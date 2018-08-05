package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	chmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	resmgmt "github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"time"
)

type FabricSetup struct {
	ConfigFile      string
	orgID           string
	ChainCodeID     string
	ChannelID       string
	initialized     bool
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName         string
	OtherOrg        string
	UserName        string
	client1         chclient.ChannelClient
	client2         chclient.ChannelClient
	RMC1           resmgmt.ResourceMgmtClient
	RMC2        resmgmt.ResourceMgmtClient
	sdk             *fabsdk.FabricSDK
}


func (setup *FabricSetup) Initialize() error {

	// check if already init 
	if setup.initialized {
		return fmt.Errorf("sdk already initialized")
	}

	// init from config 
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return fmt.Errorf("failed to create sdk: %v", err)
	}
	setup.sdk = sdk

	chMgmtClient, err := setup.sdk.NewClient(fabsdk.WithUser("Admin"), fabsdk.WithOrg(setup.OrgName)).ChannelMgmt()

	if err != nil {
		return fmt.Errorf("failed to add Admin user to sdk: %v", err)
	}

	// creating a session with channel creator org
	session, err := setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName)).Session()
	if err != nil {
		return fmt.Errorf("failed to get session for %s, %s: %s", setup.OrgName, setup.OrgAdmin, err)
	}
	orgAdminUser := session

	// request to create channel and save it 
	req := chmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfig: setup.ChannelConfig, SigningIdentity: orgAdminUser}
	if err = chMgmtClient.SaveChannel(req); err != nil {
		return fmt.Errorf("failed to create channel: %v", err)
	}

	// Allow orderer to process channel creation
	time.Sleep(time.Second * 10)

	// org1 resource management client to interact with blockchain 
	setup.RMC1, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("failed to create new resource management client: %v", err)
	}

	// org1 peers join channel
	if err = setup.RMC1.JoinChannel(setup.ChannelID); err != nil {
		return fmt.Errorf("org 1 peers failed to join the channel: %v", err)
	}
	// org2 resource manager
	setup.RMC2, err = setup.sdk.NewClient(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OtherOrg)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("org 2 failed getting resource mgmt: %v", err)
	}

	// org2 peers join channel
	if err = setup.RMC2.JoinChannel(setup.ChannelID); err != nil {
		return fmt.Errorf("org 2  peers failed to join the channel: %v", err)
	}



	fmt.Println("Initialization Successful")
	setup.initialized = true
	return nil
}



func (setup *FabricSetup) InstallAndInstantiateCC() error {

	// new cc package 
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("failed to create chaincode package: %v", err)
	}

	// rmc for org1 and org2 have to install cc on all their peers so that peers can interact with channel
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Package: ccPkg}
	_, err = setup.RMC1.InstallCC(installCCReq)
	if err != nil {
		return fmt.Errorf("failed to install cc to org1 peers %v", err)
	}
	_, err = setup.RMC2.InstallCC(installCCReq)
	if err != nil {
		return fmt.Errorf("failed to install cc to org2 peers %v", err)
	}
	// Endorse the transaction if the transaction have been signed by a member from the org "org1.hf.baltransfer.io" and "org2.hf.baltransfer.io"
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.hf.baltransfer.io"})

	// instantiate chaincode on all peers in channel, so sufficient for only one org to do it; if need to exclude some peers from this cc init, can add RequestOption and specify stuff there
	// when changing chaincode, change the version, cause it won't update if version is not updated
	err = setup.RMC1.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "1.0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
	if err != nil {
		return fmt.Errorf("failed to init cc on org1 peers %v", err)
	}
	// Channel client is used to query and execute transactions
	setup.client1, err = setup.sdk.NewClient(fabsdk.WithUser(setup.UserName), fabsdk.WithOrg(setup.OrgName)).Channel(setup.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to create new channel client for org 1: %v", err)
	}
	setup.client2, err = setup.sdk.NewClient(fabsdk.WithUser(setup.UserName), fabsdk.WithOrg(setup.OtherOrg)).Channel(setup.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to create new channel client for org 2: %v", err)
	}
	fmt.Println("Chaincode Installation & Instantiation Successful")
	return nil
}
