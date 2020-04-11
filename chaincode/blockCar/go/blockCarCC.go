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

//init;
//args: 车牌，所有人，信息id，名字，上锁，命令执行人，
func (t *BlockCarCC) initCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	carNum := args[0]

	//判断 carNum 是否存在
	carNumAsBytes, err := stub.GetState(carNum)

	if err != nil {
		return shim.Error(err.Error())
	}

	if carNumAsBytes != nil {
		return shim.Error("carCd 已经存在！")
	}

	owner := args[1]
	infoId := args[2]
	name := args[3]
	lock, err := strconv.ParseBool(args[4])
	if err != nil {
		return shim.Error("参数5必须是true or false！")
	}
	commander := args[5]
	objectType := "car"
	sensor := &def.Sensor{ObjectType: "sensor", Velocity: 0.0, Temperature: 28.0, FaultCode: "null"}
	info := &def.Infomation{ObjectType: "carInfo", ID: infoId, Name: name, Owner: owner, CarNumber: carNum}

	car := &def.Car{objectType, lock, commander, *info, *sensor}

	carJsonAsBytes, err := json.Marshal(car)

	err = stub.PutState(carNum, carJsonAsBytes)
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

//根据车牌查询car;
//args: 车牌
func (t *BlockCarCC) readCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	carNum := args[0]

	carAsBytes, err := stub.GetState(carNum)

	if err != nil {
		return shim.Error(err.Error())
	} else if carAsBytes == nil {
		return shim.Error("car 信息不存在！")
	}

	return shim.Success(carAsBytes)
}

//锁;
//args: 车牌, 上锁/开锁 （true/false）, 命令人
func (t *BlockCarCC) lockCar(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	carNum := args[0]

	carStuAsBytes, err := stub.GetState(carNum)
	if err != nil {
		return shim.Error(err.Error())
	} else if carStuAsBytes == nil {
		return shim.Error("car 不存在！")
	}

	control, err := strconv.ParseBool(args[1])
	if err != nil {
		return shim.Error("参数2 为true/false！")
	}

	commander := args[2]

	carStu := def.Car{}
	err = json.Unmarshal(carStuAsBytes, carStu)
	if err != nil {
		return shim.Error(err.Error())
	}

	carStu.Lock = control
	carStu.Commander = commander

	carStuJsonAsBytes, err := json.Marshal(carStu)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(carNum, carStuJsonAsBytes)
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
