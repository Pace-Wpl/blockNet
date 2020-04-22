package def

type CarDy struct {
	ObjectType  string  `json:"objectType"`
	Lock        bool    `json:"isLock"`      // 是否上锁
	Commander   string  `json:"commander"`   // 命令执行者
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

type CarInfomation struct {
	ObjectType string `json:"objectType"`
	Name       string `json:"name"`
	CarNumber  string `json:"carNumber"` // 车牌号
	ID         string `json:"infoId"`    // 车辆信息id
	Owner      string `json:"owner"`     // 所有人
	Type       string `json:"carType"`   // 车型号
	Colour     string `json:"carColour"` // 车颜色
}

type Regulations struct { // 违章
	ObjectType string `json:"objectType"`
	Id         string `json:"regulationsId"` // 违章id
	CarNumber  string `json:"carNumber"`     // 车牌号
	Place      string `json:"place"`         // 违章地点
	Type       string `json:"type"`          // 违章类型
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
