package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
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
//args: json CarInfomation
func (t *BlockCarCC) initCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	car := &def.CarInfomation{}
	if err := json.Unmarshal([]byte(args[0]), car); err != nil {
		return shim.Error(err.Error())
	}

	//判断 carNum 是否存在
	carInfoAsBytes, err := stub.GetState(car.ID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if carInfoAsBytes != nil {
		return shim.Error("car  已经存在！")
	}

	car.ObjectType = "carInfomation"
	carJsonAsBytes, err := json.Marshal(car)
	if err != nil {
		return shim.Error(err.Error())
	}

	carDy := &def.CarDy{ObjectType: "carDy", Lock: true, Commander: "", Velocity: 0, Temperature: 28.0, FaultCode: ""}
	carDyJSON, err := json.Marshal(carDy)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(car.ID, carJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(car.CarNumber, carDyJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent(args[1], []byte{}) //set event init
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//putCarSt 更新汽车静态信息;
//args: json carimfomation
func (t *BlockCarCC) putCarSt(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	car := &def.CarInfomation{}
	if err := json.Unmarshal([]byte(args[0]), car); err != nil {
		return shim.Error(err.Error())
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
		return shim.Error(err.Error())
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
		return shim.Error(err.Error())
	}

	//判断 carNum 是否存在
	carDyAsBytes, err := stub.GetState(carDyReq.CarNumber)
	if err != nil {
		return shim.Error(err.Error())
	}
	if carDyAsBytes == nil {
		return shim.Error("car  不存在！")
	}

	CarDy := &def.CarDy{ObjectType: "carDy", Lock: carDyReq.Lock, Commander: carDyReq.Commander, Velocity: carDyReq.Velocity, Temperature: carDyReq.Temperature, FaultCode: carDyReq.FaultCode}
	CarDyJSON, err := json.Marshal(CarDy)
	if err != nil {
		return shim.Error(err.Error())
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
//args: 所有者
func (t *BlockCarCC) queryCarByOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	owner := args[0]
	queryStr := fmt.Sprintf("{\"selector\":{\"owner\":\"%s\"}}", owner)

	resultIterator, err := stub.GetQueryResult(queryStr)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	isWrite := false

	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if isWrite == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"key\":")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString(",\"record\": ")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		isWrite = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

//根据信息id查询car;
//args: 车牌
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
	if err = json.Unmarshal(carDyAsBytes, carDy); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(resp)
}

//锁;
//args:json carDyReq
func (t *BlockCarCC) lockCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	t.putCarDy(stub, args)
	carDyReq := &def.CarDyReq{}
	if err := json.Unmarshal([]byte(args[0]), carDyReq); err != nil {
		return shim.Error(err.Error())
	}

	err := stub.SetEvent(carDyReq.CarNumber, []byte("lock")) //set event lock
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

	var buffer bytes.Buffer

	buffer.WriteString("[")

	isWrite := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()

		if err != nil {
			return shim.Error(err.Error())
		}

		if isWrite == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{ \"TxId\":")
		buffer.WriteString(queryResponse.TxId)

		buffer.WriteString(",\"Timestamp\": ")
		buffer.WriteString(time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String())

		buffer.WriteString(",\"Value\": ")
		buffer.WriteString(string(queryResponse.Value))

		buffer.WriteString(",\"IsDelete\": ")
		buffer.WriteString(strconv.FormatBool(queryResponse.IsDelete))
		buffer.WriteString("}")

		isWrite = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

//删除车辆信息;
//args: 车牌
func (t *BlockCarCC) deleteCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	carNum := args[0]

	carStuAsBytes, err := stub.GetState(carNum)

	if err != nil {
		return shim.Error(err.Error())
	} else if carStuAsBytes == nil {
		return shim.Error("car not exist!")
	}

	err = stub.DelState(carNum)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func faultHandle(stub shim.ChaincodeStubInterface, carDy *def.CarDyReq) error {
	switch carDy.FaultCode {
	case "111":
		err := stub.SetEvent(carDy.CarNumber, []byte("unknown error，fault code:"+carDy.FaultCode))
		if err != nil {
			return err
		}
	}
	return nil
}
