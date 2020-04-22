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
	_, err := ServiceClient.Client.Execute(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return err
	}
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
func InitCar(car *def.CarInit) error {

	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	log.Println(car)
	c, err := json.Marshal(car)
	if err != nil {
		return err
	}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "initCar", Args: [][]byte{c, []byte(eventID)}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return err
	}

	return nil
}

func GetCar(carNum string) (*def.CarResp, error) {
	car := &def.CarResp{}

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
