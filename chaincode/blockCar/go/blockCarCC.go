package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blockNet/chaincode/blockCar/go/def"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

func (t *BlockCarCC) testChaincode(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	value := "test chaincode!"
	err := stub.PutState(args[0], []byte(value))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("test sucess!"))
}

func (t *BlockCarCC) testChaincodeQ(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	res, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(res))
}

//init;
//args: json CarInfomation, json CarDy, event id
func (t *BlockCarCC) initCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	carInfo := &def.CarInfomation{}
	if err := json.Unmarshal([]byte(args[0]), carInfo); err != nil {
		return shim.Error(def.ErrorBadRequest)
	}

	//判断 carid 是否存在
	carInfoAsBytes, err := stub.GetState(carInfo.ID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if carInfoAsBytes != nil {
		return shim.Error("car  已经存在！")
	}

	carInfo.ObjectType = "carInfomation"
	carJsonAsBytes, err := json.Marshal(carInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(carInfo.ID, carJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(carInfo.CarNumber, []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[2], []byte{}) //set event init
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//updataCar 更新汽车静态信息;
//args: json carimfomation, event id
func (t *BlockCarCC) updataCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	car := &def.CarInfomation{}
	if err := json.Unmarshal([]byte(args[0]), car); err != nil {
		return shim.Error(def.ErrorBadRequest)
	}

	//判断 carNum 是否存在
	carInfoAsBytes, err := stub.GetState(car.ID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if carInfoAsBytes == nil {
		return shim.Error("car  不存在！")
	}

	car.ObjectType = "carInfomation"
	carJsonAsBytes, err := json.Marshal(car)
	if err != nil {
		return shim.Error(def.ErrorInternalFaults)
	}

	err = stub.PutState(car.ID, carJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[1], []byte{}) //set event
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//putCarDy 更新汽车动态信息;
//args: json carDyReq
func (t *BlockCarCC) putCarDy(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	carDyReq := &def.CarDyReq{}
	if err := json.Unmarshal([]byte(args[0]), carDyReq); err != nil {
		return shim.Error(def.ErrorBadRequest)
	}

	//判断 carNum 是否存在
	// carDyAsBytes, err := stub.GetState(carDyReq.CarNumber)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// if carDyAsBytes == nil {
	// 	return shim.Error("car  不存在！")
	// }

	CarDy := &def.CarDy{ObjectType: "carDy", Lock: carDyReq.Lock, Commander: carDyReq.Commander, Velocity: carDyReq.Velocity, Temperature: carDyReq.Temperature, FaultCode: carDyReq.FaultCode}
	CarDyJSON, err := json.Marshal(CarDy)
	if err != nil {
		return shim.Error(def.ErrorInternalFaults)
	}

	err = stub.PutState(carDyReq.CarNumber, CarDyJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = faultHandle(stub, carDyReq)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//通过所有者查询car;
//args: owner
func (t *BlockCarCC) queryCarByOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	owner := args[0]
	queryStr := fmt.Sprintf("{\"selector\":{\"owner\":\"%s\"}}", owner)

	resultIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var carItem []def.CarInfomation
	var car def.CarInfomation
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if err := json.Unmarshal(queryResponse.Value, &car); err != nil {
			return shim.Error(err.Error())
		} else {
			carItem = append(carItem, car)
		}
	}

	resp := &def.OwenrCarItem{Item: carItem}
	carItemAsBytes, err := json.Marshal(resp)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(carItemAsBytes)
}

//根据信息id查询car;
//args: 信息id
func (t *BlockCarCC) readCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	carId := args[0]

	carInfoAsBytes, err := stub.GetState(carId)

	if err != nil {
		return shim.Error(err.Error())
	} else if carInfoAsBytes == nil {
		return shim.Error("car 信息不存在！")
	}

	carInfo := &def.CarInfomation{}
	if err = json.Unmarshal(carInfoAsBytes, carInfo); err != nil {
		return shim.Error(err.Error())
	}

	carDyAsBytes, err := stub.GetState(carInfo.CarNumber)
	if err != nil {
		return shim.Error(err.Error())
	}

	carDy := &def.CarDy{}
	if err = json.Unmarshal(carDyAsBytes, carDy); err != nil {
		return shim.Error(err.Error())
	}

	rulst := &def.CarResp{
		Name: carInfo.Name, CarNumber: carInfo.CarNumber, ID: carInfo.ID, Owner: carInfo.Owner,
		Type: carInfo.Type, Colour: carInfo.Colour, Lock: carDy.Lock, Commander: carDy.Commander,
		Velocity: carDy.Velocity, Temperature: carDy.Temperature, FaultCode: carDy.FaultCode,
	}

	resp, err := json.Marshal(rulst)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(resp)
}

//根据carnum查询汽车动态信息;
//args:carnum
func (t *BlockCarCC) carState(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	carNum := args[0]

	carDyAsBytes, err := stub.GetState(carNum)
	if err != nil {
		return shim.Error(err.Error())
	} else if carDyAsBytes == nil {
		return shim.Error("car 信息不存在！")
	}

	carDy := &def.CarDy{}
	if err = json.Unmarshal(carDyAsBytes, carDy); err != nil {
		return shim.Error(err.Error())
	}

	resp, err := json.Marshal(carDy)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(resp)
}

//上锁;
//args:json carDyReq
func (t *BlockCarCC) lockCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	carNum, err := putCar(stub, args)

	err = stub.SetEvent(carNum+",islock", []byte("lock")) //set event lock
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//开锁;
//args:json carDyReq
func (t *BlockCarCC) unLockCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	carNum, err := putCar(stub, args)

	err = stub.SetEvent(carNum+",islock", []byte("unlock")) //set event lock
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//查询历史修改记录;
//args: 车牌,
func (t *BlockCarCC) getHistoryForCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	carNum := args[0]

	resultIterator, err := stub.GetHistoryForKey(carNum)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var hsItem []def.CarHistry
	var cDy def.CarDy
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var hItem def.CarHistry
		hItem.Txid = queryResponse.TxId
		hItem.IsDelete = queryResponse.IsDelete
		hItem.Timestamp = time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String()

		if err = json.Unmarshal(queryResponse.Value, &cDy); err != nil {
			return shim.Error(err.Error())
		}
		if queryResponse.Value == nil {
			var empty def.CarDy
			hItem.CarDy = empty
		} else {
			hItem.CarDy = cDy
		}

		hsItem = append(hsItem, hItem)

	}

	hstt := def.HistryItem{Item: hsItem}
	hsAsJSON, err := json.Marshal(hstt)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(hsAsJSON)
}

//删除车辆信息;
//args: 信息id,event id
func (t *BlockCarCC) deleteCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	carId := args[0]

	carStuAsBytes, err := stub.GetState(carId)
	if err != nil {
		return shim.Error(err.Error())
	} else if carStuAsBytes == nil {
		return shim.Error("car not exist!")
	}

	carInfo := &def.CarInfomation{}
	if err = json.Unmarshal(carStuAsBytes, carInfo); err != nil {
		return shim.Error(err.Error())
	}

	err = stub.DelState(carInfo.CarNumber) // 删除动态信息
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.DelState(carId) // 删除静态信息
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[1], []byte{}) //set event init
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//错误码处理
func faultHandle(stub shim.ChaincodeStubInterface, carDy *def.CarDyReq) error {
	switch carDy.FaultCode {
	case "":
		return nil
	default:
		err := stub.SetEvent(carDy.CarNumber+",faultcode", []byte("unknown error，fault code:"+carDy.FaultCode))
		if err != nil {
			return err
		}
	}
	return nil
}

//更新汽车动态信息
func putCar(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	carDyReq := &def.CarDyReq{}
	if err := json.Unmarshal([]byte(args[0]), carDyReq); err != nil {
		return "", err
	}

	//判断 carNum 是否存在
	// carDyAsBytes, err := stub.GetState(carDyReq.CarNumber)
	// if err != nil {
	// 	return err
	// }
	// if carDyAsBytes == nil {
	// 	return shim.Error("car  不存在！")
	// }

	CarDy := &def.CarDy{ObjectType: "carDy", Lock: carDyReq.Lock, Commander: carDyReq.Commander, Velocity: carDyReq.Velocity, Temperature: carDyReq.Temperature, FaultCode: carDyReq.FaultCode}
	CarDyJSON, err := json.Marshal(CarDy)
	if err != nil {
		return "", err
	}

	err = stub.PutState(carDyReq.CarNumber, CarDyJSON)
	if err != nil {
		return "", err
	}

	err = faultHandle(stub, carDyReq)
	if err != nil {
		return "", err
	}

	return carDyReq.CarNumber, nil
}
