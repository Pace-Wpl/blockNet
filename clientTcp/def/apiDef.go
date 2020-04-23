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
	LOCK_EVENT       = ",islock"
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

//request
type CarInit struct {
	Name        string  `json:"name"`
	CarNumber   string  `json:"carNumber"`   // 车牌号
	ID          string  `json:"infoId"`      // 车辆信息id
	Owner       string  `json:"owner"`       // 所有人
	Type        string  `json:"carType"`     // 车型号
	Colour      string  `json:"carColour"`   // 车颜色
	Lock        bool    `json:"isLock"`      // 是否上锁
	Commander   string  `json:"commander"`   // 命令执行者
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

type CarDyReq struct {
	CarNumber   string  `json:"carNum"`      // 车牌号
	Lock        bool    `json:"isLock"`      // 是否上锁
	Commander   string  `json:"commander"`   // 命令执行者
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

//response
type CarResp struct {
	Name        string  `json:"name"`
	CarNumber   string  `json:"carNumber"`   // 车牌号
	ID          string  `json:"infoId"`      // 车辆信息id
	Owner       string  `json:"owner"`       // 所有人
	Type        string  `json:"carType"`     // 车型号
	Colour      string  `json:"carColour"`   // 车颜色
	Lock        bool    `json:"isLock"`      // 是否上锁
	Commander   string  `json:"commander"`   // 命令执行者
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

//response
