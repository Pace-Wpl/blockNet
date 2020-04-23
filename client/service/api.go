package service

import (
	"encoding/json"
	"log"

	"github.com/blockNet/client/def"
	"github.com/blockNet/client/utils"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func TestChaincod(key string) error {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "test", Args: [][]byte{[]byte(key)}}
	log.Printf("execute key:%s\n", key)
	resp, err := ServiceClient.Client.Execute(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return err
	}
	log.Println(string(resp.Payload))
	return nil
}

func TestChaincodQ(key string) error {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "testQ", Args: [][]byte{[]byte(key)}}
	log.Printf("query key:%s\n", key)
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return err
	}
	log.Println(string(resp.Payload))
	return nil
}

//initCar;
//args: carinit struct
func InitCar(car *def.CarInit) (string, error) {

	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	carDy := &def.CarDy{ObjectType: "carDy", Lock: car.Lock, Commander: car.Commander,
		Velocity: car.Velocity, Temperature: car.Temperature, FaultCode: car.FaultCode,
	}
	carInfo := &def.CarInfomation{ObjectType: "carInfomation", Name: car.Name, CarNumber: car.CarNumber,
		ID: car.ID, Owner: car.Owner, Type: car.Type, Colour: car.Colour,
	}

	carDyAsJSON, err := json.Marshal(carDy)
	if err != nil {
		return "", err
	}

	carInfoAsJSON, err := json.Marshal(carInfo)
	if err != nil {
		return "", err
	}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "initCar", Args: [][]byte{carInfoAsJSON, carDyAsJSON, []byte(eventID)}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return "", err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func GetCar(carId string) (*def.CarInit, error) {
	carInfo := &def.CarInit{}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "readCar", Args: [][]byte{[]byte(carId)}}
	respone, err := ServiceClient.Client.Query(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return carInfo, err
	}

	if err = json.Unmarshal(respone.Payload, carInfo); err != nil {
		log.Printf("unmarshal error :%s\n", err)
		return carInfo, err
	}

	return carInfo, nil
}

func PutCarDy(carDy *def.CarDyReq) error {
	carJSON, err := json.Marshal(carDy)
	if err != nil {
		return err
	}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "putCar", Args: [][]byte{carJSON}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return err
	}

	return nil
}

//test
func LockCar(carDy *def.CarDyReq) (string, error) {
	carJSON, err := json.Marshal(carDy)
	if err != nil {
		return "", err
	}

	eventID := carDy.CarNumber + def.LOCK_EVENT
	log.Printf("eventID:%s\n", eventID)
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "lock", Args: [][]byte{carJSON}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return "", err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func QueryCarByOwner(owenr string) (*def.OwenrCarItem, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "queryCarByOwner", Args: [][]byte{[]byte(owenr)}}
	resp, err := ServiceClient.Client.Query(req)

	item := &def.OwenrCarItem{}

	if err = json.Unmarshal(resp.Payload, item); err != nil {
		return item, err
	}

	return item, nil
}

func QueryHistoryForCar(carNum string) (*def.HistryItem, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "getHistoryForCar", Args: [][]byte{[]byte(carNum)}}
	resp, err := ServiceClient.Client.Query(req)

	item := &def.HistryItem{}

	if err = json.Unmarshal(resp.Payload, item); err != nil {
		return item, err
	}

	return item, nil
}
