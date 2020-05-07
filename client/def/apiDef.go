package def

const (
	CONFIG_FILE     = "./config.yaml"
	INITIALIZED     = false
	CHAINCODE_ID    = "bcc"
	USER_NAME       = "User1"
	ORG_NAME        = "Org1"
	CHAINNEL_ID     = "mychannel"
	LOCK_EVENT      = ",islock"
	FAULTCODE_EVENT = ",faultcode"
	RGL_EVENT_SPPED = ",speedLimit"
	ON_ROAD_EVENT   = ",onRoad"
	COLLIS_EVENT    = ",collision"
)

type InitInfo struct {
	ChannelID string //通道名称
	OrgName   string //组织名称
	UserName  string //组织用户名称
}

//request
type CarInit struct {
	Name        string `json:"name"`
	CarNumber   string `json:"carNumber"` // 车牌号
	ID          string `json:"infoId"`    // 车辆信息id
	Certificate string `json:"owner"`     // 所有人证书
	Type        string `json:"carType"`   // 车型号
	Colour      string `json:"carColour"` // 车颜色
}

type OnRoadReq struct {
	Code      string  `json:"roadCode"`  // 道路编码
	CarNum    string  `json:"carNum"`    // 车牌号
	Velocity  float32 `json:"velocity"`  // 速度
	Direction string  `json:"direction"` // 方向
	Position  int     `json:"position"`  // 坐标
}

type CarDyReq struct {
	CarNumber string `json:"carNum"` // 车牌号
	// Commander   string  `json:"commander"`   // 命令执行者
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

type RglReq struct {
	RglID string `json:"rglID"`
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
type SocketResp struct {
	Mes string `json:"message"`
}

//request
type UserReg struct {
	UserName string `json:"userName"`
	PassWord string `json:"passWord"`
}
