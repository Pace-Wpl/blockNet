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

type OwenrCarItem struct {
	Item []CarInfomation `json:"item"`
}

//历史信息
type CarHistry struct {
	Txid      string `json:"Txid"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"IsDelete"`
	CarDy     CarDy  `json:"carDy"`
}
type HistryItem struct {
	Item []CarHistry `json:"item"`
}

type CarRGLItem struct {
	Item []RegulationsInfo `json:"item"`
}

//request
type CheckRGL struct {
	CarNumber string  `json:"carNumber"` // 车牌号
	Road      string  `json:"road_code"` // 地点
	Velocity  float32 `json:"velocity"`  //速度
}

//道路信息
type RoadInformation struct {
	ObjectType string  `json:"objectType"`
	Code       string  `json:"roadCode"`   // 道路编码
	Name       string  `json:"roadName"`   // 道路名字
	Coordinate string  `json:"coordinate"` // 坐标
	Type       string  `json:"roadType"`   // 道路类型
	Limit      float32 `json:"speedLimit"` // 限速
	Tag        int     `json:"roadTag"`    // 标签
}

// 违章信息
type RegulationsInfo struct {
	ObjectType string `json:"objectType"`
	ID         string `json:"regulationsId"` // 违章id
	CarNumber  string `json:"carNumber"`     // 车牌号
	Road       string `json:"road_code"`     // 违章地点
	Type       string `json:"type"`          // 违章类型
	Mes        string `json:"message"`       // 违章信息
}
