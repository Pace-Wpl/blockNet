package def

//汽车动态信息
type CarDy struct {
	ObjectType  string  `json:"objectType"`
	Certificate string  `json:"certificate"` // 命令执行者
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

//汽车静态信息
type CarInfomation struct {
	ObjectType  string `json:"objectType"`
	Name        string `json:"name"`
	CarNumber   string `json:"carNumber"`   // 车牌号
	ID          string `json:"infoId"`      // 车辆信息id
	Certificate string `json:"certificate"` // 所有人
	Type        string `json:"carType"`     // 车型号
	Colour      string `json:"carColour"`   // 车颜色
}

//历史信息
type CarHistry struct {
	Txid      string `json:"Txid"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"IsDelete"`
	CarDy     CarDy  `json:"carDy"`
}

//Lock
type Lock struct {
	ObjectType  string `json:"objectType"`
	Lc          bool   `json:"isLock"`      //锁
	Certificate string `json:"certificate"` //所有人
}

//request
type CarDyReq struct {
	CarNumber   string  `json:"carNum"`      // 车牌号
	Certificate string  `json:"certificate"` // 命令执行者
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

type LockCar struct {
	CarNum      string `json:"carNum"`      //车牌
	Lc          bool   `json:"isLock"`      // 锁
	Certificate string `json:"certificate"` // 证书
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
	Certificate string  `json:"certificate"` // 命令执行者
	Velocity    float32 `json:"velocity"`    //速度
	Temperature float32 `json:"temperature"` //温度
	FaultCode   string  `json:"faultCode"`   // 故障码
}

type OwenrCarItem struct {
	Item []CarInfomation `json:"item"`
}

type HistryItem struct {
	Item []CarHistry `json:"item"`
}
