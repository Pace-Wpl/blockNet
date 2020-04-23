package def

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

const (
	CONFIG_FILE      = "/home/pace/go/src/github.com/blockNet/client/config.yaml"
	INITIALIZED      = false
	CHAINCODE_ID     = "bcc"
	USER_NAME        = "User1"
	ORG_NAME         = "Org1"
	CHAINNEL_ID      = "mychannel"
	ORG_ADMIN        = "Admin"
	ORDERER_ORG_NAME = "orderer1.blocknet.com"
	CHAINCODE_PATH   = "github.com/blockNet/chaincode/test"
	CHANNEL_CONFIG   = "/home/pace/go/src/github.com/blockNet/channel-artifacts/channel1.tx"
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

type Car struct {
	ObjectType string     `json:"objectType"`
	Lock       bool       `json:"isLock"`    // 是否上锁
	Commander  string     `json:"commander"` // 命令执行者
	Info       Infomation `json:"info"`      // 基本信息
	Sensor     Sensor     `json:"sensor"`    // 传感器
	Owner      string     `json:"owner"`     // 所有人
}

type Sensor struct {
	ObjectType  string  `json:"objectType"`
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

type Infomation struct {
	ObjectType string `json:"objectType"`
	ID         string `json:infoID` // 车辆信息id
	Name       string `json:"name"`
	CarNumber  string `json:"carNumber"` // 车牌号
}

//request
type CarInit struct {
	CarNumber string `json:"carNumber"` // 车牌号
	Owner     string `json:"owner"`     // 所有人
	ID        string `json:infoID`      // 车辆信息id
	Name      string `json:"name"`
	Lock      string `json:"isLock"`    // 是否上锁
	Commander string `json:"commander"` // 命令执行者
}
