package sdkinit

import (
	"fmt"

	"github.com/blockNet/client/def"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

var ServiceClient ServiceSetup

func init() {
	sdk, err := setupSDK(def.CONFIG_FILE, def.INITIALIZED)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer sdk.Close()

	info := &def.InitInfo{ChannelID: def.CHAINNEL_ID, UserName: def.USER_NAME, OrgName: def.ORG_NAME}

	channelClient, err := getChannelClient(sdk, info)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	ServiceClient = ServiceSetup{
		ChaincodeID: def.CHAINCODE_ID,
		Client:      channelClient,
	}
}
