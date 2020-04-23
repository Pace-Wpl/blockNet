package def

//汽车动态信息
type CarDy struct {
	ObjectType  string  `json:"objectType"`
	Lock        bool    `json:"isLock"`      // 是否上锁
	Commander   string  `json:"commander"`   // 命令执行者
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

//汽车静态信息
type CarInfomation struct {
	ObjectType string `json:"objectType"`
	Name       string `json:"name"`
	CarNumber  string `json:"carNumber"` // 车牌号
	ID         string `json:"infoId"`    // 车辆信息id
	Owner      string `json:"owner"`     // 所有人
	Type       string `json:"carType"`   // 车型号
	Colour     string `json:"carColour"` // 车颜色
}

//request
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
