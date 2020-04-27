package def

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

//response
type CarRGLItem struct {
	Item []RegulationsInfo `json:"item"`
}
