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
	Sdk         *fabsdk.FabricSDK
}

const ChaincodeVersion = "1.0"

var ServiceClient ServiceSetup

func init() {
	sdk, err := SetupSDK(def.CONFIG_FILE, def.INITIALIZED)
	if err != nil {
		panic(err.Error())
	}

	info := &def.InitInfo{ChannelID: def.CHAINNEL_ID, UserName: def.USER_NAME, OrgName: def.ORG_NAME}

	channelClient, err := GetChannelClient(sdk, info)
	if err != nil {
		panic(err.Error())
	}

	// err = CreateChannel(sdk, info)
	// if err != nil {
	// 	log.Printf("create channel error:%s\n", err)
	// 	panic(err.Error())
	// }
	// channelClient, err := InstallAndInstantiateCC(sdk, info)
	// if err != nil {
	// 	panic(err.Error())
	// }

	ServiceClient = ServiceSetup{
		ChaincodeID: def.CHAINCODE_ID,
		Client:      channelClient,
		Sdk:         sdk,
	}
}

//setupSdk
func SetupSDK(ConfigFile string, initialized bool) (*fabsdk.FabricSDK, error) {

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
func GetChannelClient(sdk *fabsdk.FabricSDK, info *def.InitInfo) (*channel.Client, error) {

	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))
	// fmt.Println([]byte(clientChannelContext))
	fmt.Println(info.ChannelID)
	// returns a Client instance. Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}

	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")

	return channelClient, nil
}

// func CreateChannel(sdk *fabsdk.FabricSDK, info *def.InitInfo) error {

// 	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
// 	if clientContext == nil {
// 		return fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
// 	}

// 	// New returns a resource management client instance.
// 	resMgmtClient, err := resmgmt.New(clientContext)
// 	if err != nil {
// 		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
// 	}

// 	// New creates a new Client instance
// 	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
// 	if err != nil {
// 		return fmt.Errorf("根据指定的 OrgName 创建 Org MSP 客户端实例失败: %v", err)
// 	}

// 	//  Returns: signing identity
// 	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
// 	if err != nil {
// 		return fmt.Errorf("获取指定id的签名标识失败: %v", err)
// 	}

// 	// SaveChannelRequest holds parameters for save channel request
// 	channelReq := resmgmt.SaveChannelRequest{ChannelID: info.ChannelID, ChannelConfigPath: info.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
// 	// save channel response with transaction ID
// 	_, err = resMgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
// 	if err != nil {
// 		return fmt.Errorf("创建应用通道失败: %v", err)
// 	}

// 	fmt.Println("通道已成功创建，")

// 	info.OrgResMgmt = resMgmtClient

// 	// allows for peers to join existing channel with optional custom options (specific peers, filtered peers). If peer(s) are not specified in options it will default to all peers that belong to client's MSP.
// 	err = info.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
// 	if err != nil {
// 		return fmt.Errorf("Peers加入通道失败: %v", err)
// 	}

// 	fmt.Println("peers 已成功加入通道.")
// 	return nil
// }

// func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, info *def.InitInfo) (*channel.Client, error) {
// 	fmt.Println("开始安装链码......")
// 	// creates new go lang chaincode package
// 	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("创建链码包失败: %v", err)
// 	}

// 	// contains install chaincode request parameters
// 	installCCReq := resmgmt.InstallCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: ChaincodeVersion, Package: ccPkg}
// 	// allows administrators to install chaincode onto the filesystem of a peer
// 	_, err = info.OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
// 	if err != nil {
// 		return nil, fmt.Errorf("安装链码失败: %v", err)
// 	}

// 	fmt.Println("指定的链码安装成功")
// 	fmt.Println("开始实例化链码......")

// 	//  returns a policy that requires one valid
// 	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.blocknet.com"})

// 	instantiateCCReq := resmgmt.InstantiateCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: ChaincodeVersion, Args: [][]byte{[]byte("init")}, Policy: ccPolicy}
// 	// instantiates chaincode with optional custom options (specific peers, filtered peers, timeout). If peer(s) are not specified
// 	_, err = info.OrgResMgmt.InstantiateCC(info.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
// 	if err != nil {
// 		return nil, fmt.Errorf("实例化链码失败: %v", err)
// 	}

// 	fmt.Println("链码实例化成功")

// 	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))
// 	// returns a Client instance. Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
// 	channelClient, err := channel.New(clientChannelContext)
// 	if err != nil {
// 		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
// 	}

// 	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")

// 	return channelClient, nil
// }
