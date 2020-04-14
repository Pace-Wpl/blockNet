package def

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

const (
	CONFIG_FILE  = "../config.yaml"
	INITIALIZED  = false
	CHAINCODE_ID = "mycc"
	USER_NAME    = "pace"
	ORG_NAME     = "Org1"
	CHAINNEL_ID  = "mychannel1"
)

type InitInfo struct {
	ChannelID      string          //通道名称
	ChannelConfig  string          //通道交易配置文件所在路径
	OrgName        string          //组织名称
	OrgAdmin       string          //组织管理员名称
	OrdererOrgName string          //Orderer名称
	OrgResMgmt     *resmgmt.Client //资源管理端实例

	ChaincodeID     string //链码ID（即链码名称）
	ChaincodeGoPath string //系统GOPATH路径
	ChaincodePath   string //链 码源代码所在路径
	UserName        string //组织用户名称
}
