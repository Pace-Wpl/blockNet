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
	switch fn {
	case "initCar":
		return t.initCar(stub, args)
	case "readCar":
		return t.readCar(stub, args)
	case "lock":
		return t.lockCar(stub, args)
	case "unlock":
		return t.unLockCar(stub, args)
	case "deleteCar":
		return t.deleteCar(stub, args)
	case "putCar":
		return t.putCarDy(stub, args)
	case "updataCar":
		return t.updataCar(stub, args)
	case "queryCarByOwner":
		return t.queryCarByOwner(stub, args) //debug
	case "getHistoryForCar":
		return t.getHistoryForCar(stub, args) //debug
	case "test":
		return t.testChaincode(stub, args)
	case "testQ":
		return t.testChaincodeQ(stub, args)
	default:
	}
	return shim.Error("没有相应的方法！")
}

func main() {

	err := shim.Start(new(BlockCarCC))
	if err != nil {
		fmt.Println("chaincode start error!")
	}
}
