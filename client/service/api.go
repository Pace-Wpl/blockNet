package service

import (
	"encoding/json"
	"log"

	"github.com/blockNet/client/def"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func InitCar(car *def.CarInit) error {

	// eventID := utils.NewUUID()
	// reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	// defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "initCar", Args: [][]byte{[]byte(car.CarNumber), []byte(car.Owner), []byte(car.ID), []byte(car.Name), []byte(car.Lock), []byte(car.Commander)}}
	_, err := ServiceClient.Client.Execute(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return err
	}

	return nil
}

func GetCar(carNum string) (*def.Car, error) {
	car := &def.Car{}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "readCar", Args: [][]byte{[]byte(carNum)}}
	respone, err := ServiceClient.Client.Query(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return car, err
	}

	if err = json.Unmarshal(respone.Payload, car); err != nil {
		log.Printf("unmarshal error :%s\n", err)
		return car, err
	}

	return car, nil
}
