package def

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

//request
type CheckRGL struct {
	CarNumber string  `json:"carNumber"` // 车牌号
	Road      string  `json:"road_code"` // 地点
	Velocity  float32 `json:"velocity"`  //速度
}
