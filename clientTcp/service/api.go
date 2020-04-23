package service

import (
	"encoding/json"
	"log"

	"github.com/blockNet/clientTcp/def"
	"github.com/blockNet/clientTcp/utils"
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

func GetCar(carId string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "readCar", Args: [][]byte{[]byte(carId)}}
	respone, err := ServiceClient.Client.Query(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return []byte{}, err
	}

	return respone.Payload, nil
}

func GetState(carNum string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "carState", Args: [][]byte{[]byte(carNum)}}
	respone, err := ServiceClient.Client.Query(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return []byte{}, err
	}

	return respone.Payload, nil
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

func UnLockCar(carDy *def.CarDyReq) (string, error) {
	carJSON, err := json.Marshal(carDy)
	if err != nil {
		return "", err
	}

	eventID := carDy.CarNumber + def.LOCK_EVENT
	log.Printf("eventID:%s\n", eventID)
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "unlock", Args: [][]byte{carJSON}}
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

func QueryCarByOwner(owenr string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "queryCarByOwner", Args: [][]byte{[]byte(owenr)}}
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}

	// item := &def.OwenrCarItem{}

	// if err = json.Unmarshal(resp.Payload, item); err != nil {
	// 	return item, err
	// }
	return resp.Payload, nil
}

func QueryHistoryForCar(carNum string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "getHistoryForCar", Args: [][]byte{[]byte(carNum)}}
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}
	// item := &def.HistryItem{}

	// if err = json.Unmarshal(resp.Payload, item); err != nil {
	// 	return item, err
	// }

	return resp.Payload, nil
}

func DeleteCar(carId string) ([]byte, error) {
	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "deleteCar", Args: [][]byte{[]byte(carId), []byte(eventID)}}
	_, err := ServiceClient.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return []byte{}, err
	}

	return []byte(resp), nil
}

func CheckRGL(cGRL *def.CheckRGL) ([]byte, error) {
	cGRLJSON, err := json.Marshal(cGRL)
	if err != nil {
		return []byte{}, err
	}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "checkRGL", Args: [][]byte{cGRLJSON}}
	resp, err := ServiceClient.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}

	return resp.Payload, nil
}

func UpdataCar(carInfo *def.CarInfomation) ([]byte, error) {
	cInfoJSON, err := json.Marshal(carInfo)
	if err != nil {
		return []byte{}, err
	}

	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "updataCar", Args: [][]byte{cInfoJSON, []byte(eventID)}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return []byte{}, err
	}

	return []byte(resp), nil
}

func UpdataRoad(roadInfo *def.RoadInformation) ([]byte, error) {
	rInfoJSON, err := json.Marshal(roadInfo)
	if err != nil {
		return []byte{}, err
	}

	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "updataRoad", Args: [][]byte{rInfoJSON, []byte(eventID)}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return []byte{}, err
	}

	return []byte(resp), nil
}

func GetRoad(roadCode string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "readRoad", Args: [][]byte{[]byte(roadCode)}}
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}

	return resp.Payload, nil
}

func DeleteRoad(roadCode string) ([]byte, error) {
	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "deleteRoad", Args: [][]byte{[]byte(roadCode), []byte(eventID)}}
	_, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return []byte{}, err
	}

	return []byte(resp), nil
}

func OnRoad(roadCode string, carNum string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "onRoad", Args: [][]byte{[]byte(roadCode), []byte(carNum)}}
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}

	return resp.Payload, nil
}

func UpdataRGL(rgl *def.RegulationsInfo) ([]byte, error) {
	rglJSON, err := json.Marshal(rgl)
	if err != nil {
		return []byte{}, err
	}

	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "updataRGL", Args: [][]byte{[]byte(rglJSON), []byte(eventID)}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return []byte{}, err
	}

	return []byte(resp), nil
}

func GetRGL(rglId string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "readRGL", Args: [][]byte{[]byte(rglId)}}
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}

	return resp.Payload, nil
}

func DealRGL(rglId string) ([]byte, error) {
	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "deleteRGL", Args: [][]byte{[]byte(rglId), []byte(eventID)}}
	_, err := ServiceClient.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return []byte{}, err
	}

	return []byte(resp), nil
}

func GetHistoryRGL(rglId string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "getHistoryRGL", Args: [][]byte{[]byte(rglId)}}
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}

	return resp.Payload, nil
}
