package def

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
