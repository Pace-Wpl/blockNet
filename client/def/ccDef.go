package def

type LockCarReq struct {
	CarNum      string `json:"carNum"`      //车牌
	Certificate string `json:"certificate"` // 证书
}

type LockCar struct {
	ObjectType  string `json:"objectType"`
	CarNum      string `json:"carNum"`      //车牌
	Lc          bool   `json:"isLock"`      // 锁
	Certificate string `json:"certificate"` // 证书
}

//汽车动态信息
type CarDy struct {
	ObjectType string `json:"objectType"`
	// Certificate string  `json:"certificate"` // 命令执行者
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

//Lock
type Lock struct {
	ObjectType  string `json:"objectType"`
	Lc          bool   `json:"isLock"`      //锁
	Certificate string `json:"certificate"` //所有人
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

// 违章信息
type RegulationsInfo struct {
	ObjectType string `json:"objectType"`
	ID         string `json:"regulationsId"` // 违章id
	CarNumber  string `json:"carNumber"`     // 车牌号
	Road       string `json:"road_code"`     // 违章地点
	Type       string `json:"type"`          // 违章类型
	Mes        string `json:"message"`       // 违章信息
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
	Limit      float32 `json:"speedLimit"` // 限速
}

//OnRoad
type OnRoad struct {
	ObjectType string  `json:"objectType"`
	Code       string  `json:"roadCode"`  // 道路编码
	CarNum     string  `json:"carNum"`    // 车牌号
	Velocity   float32 `json:"velocity"`  // 速度
	Direction  string  `json:"direction"` // 方向
	Position   int     `json:"position"`  // 坐标
}

//requets
type RoadReq struct {
	Code      string  `json:"roadCode"`  // 道路编码
	Velocity  float32 `json:"velocity"`  // 速度
	Direction string  `json:"direction"` // 方向
	Position  int     `json:"position"`  // 坐标
}

//历史信息
type HistryRGL struct {
	Txid      string          `json:"Txid"`
	Timestamp string          `json:"timestamp"`
	IsDelete  bool            `json:"IsDelete"`
	Rgl       RegulationsInfo `json:"rgl"`
}

//response
type HistryRGLItem struct {
	Item []HistryRGL `json:"item"`
}

type User struct {
	ObjectType  string `json:"objectType"`
	Certificate string `json:"certificate"`
	UserName    string `json:"userName"`
	PassWord    string `json:"passWord"`
}
