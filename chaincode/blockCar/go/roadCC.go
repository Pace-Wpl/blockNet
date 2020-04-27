package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/blockNet/chaincode/blockCar/go/def"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

//updataRoad 更新道路静态信息;
//args: json RoadInformation, event id
func (t *BlockCarCC) updataRoad(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	road := &def.RoadInformation{}
	if err := json.Unmarshal([]byte(args[0]), road); err != nil {
		return shim.Error(def.ErrorBadRequest)
	}

	//判断  是否存在
	// roadInfoAsBytes, err := stub.GetState(road.Code)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// if roadInfoAsBytess == nil {
	// 	return shim.Error("road  不存在！")
	// }

	road.ObjectType = "roadInfomation"
	roadJsonAsBytes, err := json.Marshal(road)
	if err != nil {
		return shim.Error(def.ErrorInternalFaults)
	}

	err = stub.PutState(road.Code, roadJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[1], []byte{}) //set event
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//删除道路信息;
//args: road code,
func (t *BlockCarCC) deleteRoad(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	rc := args[0]

	rcAsBytes, err := stub.GetState(rc)
	if err != nil {
		return shim.Error(err.Error())
	} else if rcAsBytes == nil {
		return shim.Error("road not exist!")
	}

	err = stub.DelState(rc) // 删除道路信息
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[1], []byte{}) //set event init
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//根据road code查询道路信息;
//args: road code,
func (t *BlockCarCC) readRoad(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	rc := args[0]

	roadInfoAsBytes, err := stub.GetState(rc)

	if err != nil {
		return shim.Error(err.Error())
	} else if roadInfoAsBytes == nil {
		return shim.Error(" 信息不存在！")
	}

	roadInfo := &def.RoadInformation{}
	if err = json.Unmarshal(roadInfoAsBytes, roadInfo); err != nil {
		return shim.Error(err.Error())
	}

	resp, err := json.Marshal(roadInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(resp)
}

//在路上,汽车传输正在行使的道路信息 ,检测超速，追尾
//args: OnRoad
// func (t *BlockCarCC) onRoad(stub shim.ChaincodeStubInterface, args []string) peer.Response {
// 	or := &def.OnRoad{}

// 	if err := json.Unmarshal([]byte(args[0]), or); err != nil {
// 		return shim.Error(def.ErrorInternalFaults)
// 	}

// 	roadInfoAsBytes, err := stub.GetState(or.Code)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	} else if roadInfoAsBytes == nil {
// 		return shim.Error("can not discern the road!")
// 	}

// 	//追尾检测
// 	t.checkCollision(stub, or)
// 	//超速检测
// 	t.checkRGL(stub, or)

// 	return shim.Success(nil)
// }

//追尾检测
//args: json Onroad
func (t *BlockCarCC) checkCollision(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	o := &def.OnRoad{}

	if err := json.Unmarshal([]byte(args[0]), o); err != nil {
		return shim.Error(def.ErrorInternalFaults)
	}
	objType := "onRoad"
	queryStr := fmt.Sprintf("{\"selector\":{\"objectType\":\"%s\",\"roadCode\":\"%s\",\"direction\":\"%s\"}}", objType, o.Code, o.Direction)

	resultIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var or def.OnRoad
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if err := json.Unmarshal(queryResponse.Value, &or); err != nil {
			return shim.Error(err.Error())
		}
		if o.CarNum != or.CarNum {
			distance := float64(o.Position - or.Position)
			ve := float64(o.Velocity - or.Velocity)
			if checkCollisionAlg(distance, ve, 3) { //追尾检测
				err = stub.SetEvent(o.CarNum+",collision", []byte{})
				if err != nil {
					return shim.Error(err.Error())
				}
			}
		}
	}

	key := o.CarNum + ",onRoad"
	oAsJSON, err := json.Marshal(o)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(key, oAsJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//追尾检测算法，
//args:距离，速度差，反映时间
func checkCollisionAlg(distance, ve, duration float64) bool {
	if math.Abs(ve*duration) >= math.Abs(distance) {
		return true
	}
	return false
}

func (t *BlockCarCC) testCollision(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	objType := "onRoad"
	queryStr := fmt.Sprintf("{\"selector\":{\"objectType\":\"%s\",\"roadCode\":\"%s\",\"direction\":\"%s\"}}", objType, args[0], args[1])
	resultIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var or def.OnRoad
	var item []def.OnRoad
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if err := json.Unmarshal(queryResponse.Value, &or); err != nil {
			return shim.Error(err.Error())
		}

		item = append(item, or)
	}

	itemAsJSON, err := json.Marshal(item)

	return shim.Success(itemAsJSON)
}

//违法检测：超速
//args:json Onroads,
func (t *BlockCarCC) checkRGL(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	c := &def.OnRoad{}

	if err := json.Unmarshal([]byte(args[0]), c); err != nil {
		return shim.Error(def.ErrorInternalFaults)
	}

	roadAsByte, err := stub.GetState(c.Code)
	if err != nil {
		return shim.Error(err.Error())
	}

	roadInfo := &def.RoadInformation{}
	if err = json.Unmarshal(roadAsByte, roadInfo); err != nil {
		return shim.Error(err.Error())
	}

	if c.Velocity > roadInfo.Limit {
		id := func(len int) string { //5 位随即字符串
			bytes := make([]byte, len)
			for i := 0; i < len; i++ {
				c := int64(i)
				b := rand.New(rand.NewSource(time.Now().Unix()+c)).Intn(26) + 65
				bytes[i] = byte(b)
			}
			return string(bytes)
		}(5)
		rglInfo := &def.RegulationsInfo{
			ObjectType: "RegulationsInfo",
			ID:         id,
			CarNumber:  c.CarNum, Road: c.Code, Type: "speeding", Mes: "speeding in the road :" + roadInfo.Name + "! RGL id :" + id,
		}

		rglAsBytes, err := json.Marshal(rglInfo)
		if err != nil {
			return shim.Error(def.ErrorInternalFaults)
		}

		err = stub.PutState(rglInfo.ID, rglAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.SetEvent(c.CarNum+",speedLimit", []byte("speeding in the road :"+roadInfo.Name+"!"))
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

//updataRGL 更新违法信息信息;
//args: json RegulationsInfo, event id
func (t *BlockCarCC) updataRGL(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	rgl := &def.RegulationsInfo{}
	if err := json.Unmarshal([]byte(args[0]), rgl); err != nil {
		return shim.Error(def.ErrorBadRequest)
	}

	// 判断  违章是否存在
	rglInfoAsBytes, err := stub.GetState(rgl.ID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if rglInfoAsBytes != nil {
		return shim.Error("违章信息已存在！")
	}

	//判断车牌是否存在
	carInfoAsBytes, err := stub.GetState(rgl.CarNumber)
	if err != nil {
		return shim.Error(err.Error())
	}
	if carInfoAsBytes == nil {
		return shim.Error("车牌不存在，违章信息有误！")
	}

	rgl.ObjectType = "RegulationsInfo"
	rglJsonAsBytes, err := json.Marshal(rgl)
	if err != nil {
		return shim.Error(def.ErrorInternalFaults)
	}

	err = stub.PutState(rgl.ID, rglJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[1], []byte{}) //set event
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//处理消除违法信息;
//args: RGL ID, event id
func (t *BlockCarCC) dealRGL(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	rgl := args[0]

	rglAsBytes, err := stub.GetState(rgl)
	if err != nil {
		return shim.Error(err.Error())
	} else if rglAsBytes == nil {
		return shim.Error("rgl not exist!")
	}

	err = stub.DelState(rgl) // 删除信息
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[1], []byte{}) //set event init
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//根据rgl id查询违法信息;
//args: rgl id,
func (t *BlockCarCC) readRgl(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	rgl := args[0]

	rglInfoAsBytes, err := stub.GetState(rgl)
	if err != nil {
		return shim.Error(err.Error())
	} else if rglInfoAsBytes == nil {
		return shim.Error(" 信息不存在！")
	}

	rglInfo := &def.RegulationsInfo{}
	if err = json.Unmarshal(rglInfoAsBytes, rglInfo); err != nil {
		return shim.Error(err.Error())
	}

	resp, err := json.Marshal(rglInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(resp)
}

//通过carNum查询rgl;
//args: carNum,event id,
func (t *BlockCarCC) carRgl(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	carNum := args[0]
	objType := "RegulationsInfo"
	queryStr := fmt.Sprintf("{\"selector\":{\"objectType\":\"%s\",\"carNumber\":\"%s\"}}", objType, carNum)

	resultIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var rglItem []def.RegulationsInfo
	var rgl def.RegulationsInfo
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if err := json.Unmarshal(queryResponse.Value, &rgl); err != nil {
			return shim.Error(err.Error())
		} else {
			rglItem = append(rglItem, rgl)
		}
	}
	resp := &def.CarRGLItem{Item: rglItem}
	rglItemAsBytes, err := json.Marshal(resp)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(rglItemAsBytes)
}

//根据车牌查询违法信息记录
//args: rgl code  ,
func (t *BlockCarCC) getHistryRgl(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	rglC := args[0]

	resultIterator, err := stub.GetHistoryForKey(rglC)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var hsItem []def.HistryRGL
	var rgl def.RegulationsInfo
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var hItem def.HistryRGL
		hItem.Txid = queryResponse.TxId
		hItem.IsDelete = queryResponse.IsDelete
		hItem.Timestamp = time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String()

		if !hItem.IsDelete {
			if err = json.Unmarshal(queryResponse.Value, &rgl); err != nil {
				return shim.Error(err.Error())
			}
			if queryResponse.Value == nil {
				var empty def.RegulationsInfo
				hItem.Rgl = empty
			} else {
				hItem.Rgl = rgl
			}
		} else {
			var empty def.RegulationsInfo
			hItem.Rgl = empty
		}

		hsItem = append(hsItem, hItem)

	}

	hstt := def.HistryRGLItem{Item: hsItem}
	hsAsJSON, err := json.Marshal(hstt)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(hsAsJSON)
}
