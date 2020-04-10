package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type AirCdChaincode struct {
}

type airCd struct {
	ObjectType  string `json:"objectType"`
	Name        string `json:"name"`
	Temperature int    `json:"temperature"`
	Open        bool   `json:"open"`
	Commander   string `json:"commander"`
	Owner       string `json:"owner"`
}

func (t *AirCdChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *AirCdChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters()
	if fn == "initAirCd" {
		return t.initAirCd(stub, args)
	} else if fn == "queryAirCdByOwner" {
		return t.queryAirCdByOwner(stub, args)
	} else if fn == "readAirCd" {
		return t.readAirCd(stub, args)
	} else if fn == "controlAirCd" {
		return t.controlAirCd(stub, args)
	} else if fn == "getHistoryForAirCd" {
		return t.getHistoryForAirCd(stub, args)
	} else if fn == "deleteAirCd" {
		return t.deleteAirCd(stub, args)
	} else if fn == "queryAirCdByRange" {
		return t.queryAirCdByRange(stub, args)
	}
	return shim.Error("没有相应的方法！")
}

func (t *AirCdChaincode) initAirCd(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	airCdname := args[0]

	//判断 airCd 是否存在
	airCdnameAsBytes, err := stub.GetState(airCdname)

	if err != nil {
		return shim.Error(err.Error())
	}

	if airCdnameAsBytes != nil {
		return shim.Error("airCd 已经存在！")
	}

	temperature, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("参数2必须是数字！")
	}

	open, err := strconv.ParseBool(args[2])
	if err != nil {
		return shim.Error("参数3必须是true or false！")
	}

	commander := args[3]
	owner := args[4]
	objectType := "airCd"
	airCd := &airCd{objectType, airCdname, temperature, open, commander, owner}

	airCdJsonAsBytes, err := json.Marshal(airCd)

	err = stub.PutState(airCdname, airCdJsonAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *AirCdChaincode) queryAirCdByOwner(stub shim.ChaincodeStubInterface, args []string) peer.Response {

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

func (t *AirCdChaincode) readAirCd(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	airCdname := args[0]

	airCdAsBytes, err := stub.GetState(airCdname)

	if err != nil {
		return shim.Error(err.Error())
	} else if airCdAsBytes == nil {
		return shim.Error("airCd 信息不存在！")
	}

	return shim.Success(airCdAsBytes)
}

func (t *AirCdChaincode) controlAirCd(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	control := args[1]

	if control == "temperature" {
		return t.controlAirCdBytp(stub, args)
	} else if control == "open" {
		return t.controlAirCdByOp(stub, args)
	}

	return shim.Error("paramete_2 error!")
}

func (t *AirCdChaincode) controlAirCdBytp(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	airCdname := args[0]

	airCdAsBytes, err := stub.GetState(airCdname)

	if err != nil {
		return shim.Error(err.Error())
	} else if airCdAsBytes == nil {
		return shim.Error("airCd 信息不存在！")
	}

	paramete, err := strconv.Atoi(args[2])

	if err != nil {
		return shim.Error("paramete_3 must be number! ")
	}

	commander := args[3]

	airCdinfo := airCd{}
	err = json.Unmarshal(airCdAsBytes, &airCdinfo)

	if err != nil {
		return shim.Error(err.Error())
	}

	airCdinfo.Temperature = paramete
	airCdinfo.Commander = commander

	airCdJsonAsBytes, err := json.Marshal(airCdinfo)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(airCdname, airCdJsonAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *AirCdChaincode) controlAirCdByOp(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	airCdname := args[0]

	airCdAsBytes, err := stub.GetState(airCdname)

	if err != nil {
		return shim.Error(err.Error())
	} else if airCdAsBytes == nil {
		return shim.Error("aircd not exist!")
	}

	paramete, err := strconv.ParseBool(args[2])

	if err != nil {
		return shim.Error("paramete_3 must be true or false!")
	}

	commander := args[3]

	airCdinfo := airCd{}

	err = json.Unmarshal(airCdAsBytes, &airCdinfo)

	if err != nil {
		return shim.Error(err.Error())
	}

	airCdinfo.Open = paramete

	airCdinfo.Commander = commander

	airCdJsonAsBytes, err := json.Marshal(airCdinfo)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(airCdname, airCdJsonAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *AirCdChaincode) getHistoryForAirCd(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	airCdname := args[0]

	resultIterator, err := stub.GetHistoryForKey(airCdname)

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

func (t *AirCdChaincode) deleteAirCd(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	airCdname := args[0]

	airCdAsBytes, err := stub.GetState(airCdname)

	if err != nil {
		return shim.Error(err.Error())
	} else if airCdAsBytes == nil {
		return shim.Error("airCd not exist!")
	}

	err = stub.DelState(airCdname)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *AirCdChaincode) queryAirCdByRange(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	startKey := args[0]
	endKey := args[1]

	resultIterator, err := stub.GetStateByRange(startKey, endKey)

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

		buffer.WriteString("{ \"key\": ")
		buffer.WriteString(queryResponse.Key)

		buffer.WriteString(",\"record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		isWrite = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func main() {

	err := shim.Start(new(AirCdChaincode))
	if err != nil {
		fmt.Println("chaincode start error!")
	}
}
