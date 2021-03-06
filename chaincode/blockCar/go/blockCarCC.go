package main

import (
	"encoding/json"
	"errors"
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
//args: json CarInfomation, json CarDy, json LockCar, event id
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

	//判断 cer 是否存在
	exit, err := CerExit(stub, carInfo.Certificate)
	if !exit || err != nil {
		return shim.Error("cer 不存在!")
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

	key := carInfo.CarNumber + ";lock"
	err = stub.PutState(key, []byte(args[2]))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[3], []byte{}) //set event init
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//判断证书是否存在;
//args: certs
func CerExit(stub shim.ChaincodeStubInterface, cert string) (bool, error) {
	queryStr := fmt.Sprintf("{\"selector\":{\"objectType\":\"user\",\"certificate\":\"%s\"}}", cert)

	resultIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return false, err
	}
	defer resultIterator.Close()

	if resultIterator.HasNext() {
		return true, nil
	}

	return false, nil
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

	CarDy := &def.CarDy{ObjectType: "carDy", Velocity: carDyReq.Velocity, Temperature: carDyReq.Temperature, FaultCode: carDyReq.FaultCode}
	CarDyJSON, err := json.Marshal(CarDy)
	if err != nil {
		return shim.Error(def.ErrorInternalFaults)
	}

	err = stub.PutState(carDyReq.CarNumber, CarDyJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// err = faultHandle(stub, carDyReq)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	return shim.Success(nil)
}

//通过所有者查询car;
//args: owner
func (t *BlockCarCC) queryCarByOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	owner := args[0]
	objType := "carInfomation"
	queryStr := fmt.Sprintf("{\"selector\":{\"objectType\":\"%s\",\"owner\":\"%s\"}}", objType, owner)

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

	key := carInfo.CarNumber + ";lock"
	carLockJSON, err := stub.GetState(key)
	carLock := &def.Lock{}
	if err = json.Unmarshal(carLockJSON, carLock); err != nil {
		return shim.Error(err.Error())
	}

	rulst := &def.CarResp{
		Name: carInfo.Name, CarNumber: carInfo.CarNumber, ID: carInfo.ID, Owner: carInfo.Certificate,
		Type: carInfo.Type, Colour: carInfo.Colour, Lock: carLock.Lc,
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

//上/开 锁;
//args:json LockCar
func (t *BlockCarCC) lockCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	carNum, err := lockCar(stub, []byte(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(carNum+",islock", []byte("lock")) //set event lock
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//开锁;
//args:json LockCar
func (t *BlockCarCC) unLockCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	carNum, err := lockCar(stub, []byte(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	}

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

		if !hItem.IsDelete {
			if err = json.Unmarshal(queryResponse.Value, &cDy); err != nil {
				return shim.Error(err.Error())
			}
			if queryResponse.Value == nil {
				var empty def.CarDy
				hItem.CarDy = empty
			} else {
				hItem.CarDy = cDy
			}
		} else {
			var empty def.CarDy
			hItem.CarDy = empty
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
// func faultHandle(stub shim.ChaincodeStubInterface, carDy *def.CarDyReq) error {
// 	switch carDy.FaultCode {
// 	case "":
// 		return nil
// 	default:
// 		err := stub.SetEvent(carDy.CarNumber+",faultcode", []byte("unknown error，fault code:"+carDy.FaultCode))
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

//更新锁车动态信息
func lockCar(stub shim.ChaincodeStubInterface, args []byte) (string, error) {
	carDyReq := &def.LockCar{}
	if err := json.Unmarshal(args, carDyReq); err != nil {
		return "", err
	}

	key := carDyReq.CarNum + ";lock"
	carDyAsBytes, err := stub.GetState(key)
	if err != nil {
		return "", err
	} else if carDyAsBytes == nil {
		return "", errors.New("信息不存在")
	}

	carDy := &def.Lock{}
	if err = json.Unmarshal(carDyAsBytes, carDy); err != nil {
		return "", err
	}

	if carDyReq.Certificate != carDy.Certificate {
		return "", errors.New("操作失败，没有权限！")
	} else {
		carDy.Lc = carDyReq.Lc
	}

	//判断 carNum 是否存在
	// carDyAsBytes, err := stub.GetState(carDyReq.CarNumber)
	// if err != nil {
	// 	return err
	// }
	// if carDyAsBytes == nil {
	// 	return shim.Error("car  不存在！")
	// }

	CarDyJSON, err := json.Marshal(carDy)
	if err != nil {
		return "", err
	}

	err = stub.PutState(key, CarDyJSON)
	if err != nil {
		return "", err
	}

	return carDyReq.CarNum, nil
}

//
//args: carNum
func (t *BlockCarCC) getLock(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	key := args[0] + ";lock"

	LockJSON, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(LockJSON)
}

//args: json user,
func (t *BlockCarCC) register(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	user := &def.User{}
	if err := json.Unmarshal([]byte(args[0]), user); err != nil {
		return shim.Error(def.ErrorBadRequest)
	}

	//判断 usr 是否存在
	uAsBytes, err := stub.GetState(user.UserName)
	if err != nil {
		return shim.Error(err.Error())
	}
	if uAsBytes != nil {
		return shim.Error("user name had exeit！")
	}

	err = stub.PutState(user.UserName, []byte(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[1], []byte{}) //set event init
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *BlockCarCC) getUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	uRq := &def.UserReq{}
	if err := json.Unmarshal([]byte(args[0]), uRq); err != nil {
		return shim.Error(def.ErrorBadRequest)
	}

	uAsBytes, err := stub.GetState(uRq.UserName)
	if err != nil {
		return shim.Error(err.Error())
	}

	u := &def.User{}
	if err := json.Unmarshal(uAsBytes, u); err != nil {
		return shim.Error(err.Error())
	}

	if uRq.PassWord != u.PassWord {
		return shim.Error("pass word error!")
	}

	return shim.Success(uAsBytes)
}
