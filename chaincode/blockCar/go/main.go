package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type BlockCarCC struct {
}

func (t *BlockCarCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *BlockCarCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters()
	if fn == "initCar" {
		return t.initCar(stub, args)
	} else if fn == "queryCarByOwner" {
		return t.queryCarByOwner(stub, args)
	} else if fn == "readCar" {
		return t.readCar(stub, args)
	} else if fn == "lock" {
		return t.lockCar(stub, args)
	} else if fn == "getHistoryForCar" {
		return t.getHistoryForCar(stub, args)
	} else if fn == "deleteCar" {
		return t.deleteCar(stub, args)
	} else if fn == "putCar" {
		return t.putCarDy(stub, args)
	} else if fn == "test" {
		return t.testChaincode(stub, args)
	} else if fn == "testQ" {
		return t.testChaincodeQ(stub, args)
	}
	return shim.Error("没有相应的方法！")
}

func main() {

	err := shim.Start(new(BlockCarCC))
	if err != nil {
		fmt.Println("chaincode start error!")
	}
}
