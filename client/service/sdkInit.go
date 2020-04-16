package service

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/blockNet/client/def"
)

type ServiceSetup struct {
	ChaincodeID string
	Client      *channel.Client
}

var ServiceClient ServiceSetup

func init() {
	sdk, err := setupSDK(def.CONFIG_FILE, def.INITIALIZED)
	if err != nil {
		panic(err.Error())
	}
	defer sdk.Close()

	info := &def.InitInfo{ChannelID: def.CHAINNEL_ID, UserName: def.USER_NAME, OrgName: def.ORG_NAME}

	channelClient, err := getChannelClient(sdk, info)
	if err != nil {
		panic(err.Error())
	}

	ServiceClient = ServiceSetup{
		ChaincodeID: def.CHAINCODE_ID,
		Client:      channelClient,
	}
}

//setupSdk
func setupSDK(ConfigFile string, initialized bool) (*fabsdk.FabricSDK, error) {

	if initialized {
		return nil, fmt.Errorf("Fabric SDK已被实例化")
	}

	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败: %v", err)
	}

	fmt.Println("Fabric SDK初始化成功")
	return sdk, nil
}

//获取通道客户端
func getChannelClient(sdk *fabsdk.FabricSDK, info *def.InitInfo) (*channel.Client, error) {

	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))
	// returns a Client instance. Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}

	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")

	return channelClient, nil
}
