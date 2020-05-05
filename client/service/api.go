package service

import (
	"encoding/json"
	"log"
	"net"
	"sync"

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

	carDy := &def.CarDy{ObjectType: "carDy",
		Velocity: 0.0, Temperature: 28.3, FaultCode: "",
	}
	carInfo := &def.CarInfomation{ObjectType: "carInfomation", Name: car.Name, CarNumber: car.CarNumber,
		ID: car.ID, Certificate: car.Certificate, Type: car.Type, Colour: car.Colour,
	}

	lock := &def.LockCar{ObjectType: "carLock", Lc: false, Certificate: car.Certificate, CarNum: car.CarNumber}

	lockJSON, err := json.Marshal(lock)
	if err != nil {
		return "", err
	}

	carDyAsJSON, err := json.Marshal(carDy)
	if err != nil {
		return "", err
	}

	carInfoAsJSON, err := json.Marshal(carInfo)
	if err != nil {
		return "", err
	}

	log.Println(string(carInfoAsJSON))
	log.Println(string(carDyAsJSON))
	log.Println(string(lockJSON))
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "initCar", Args: [][]byte{carInfoAsJSON, carDyAsJSON, lockJSON, []byte(eventID)}}
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
func LockCar(carDy *def.LockCar) (string, error) {
	carDy.Lc = true
	carJSON, err := json.Marshal(carDy)
	if err != nil {
		return "", err
	}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "lock", Args: [][]byte{carJSON}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return "", err
	}

	return "successful!", nil
}

func UnLockCar(carDy *def.LockCar) (string, error) {
	carDy.Lc = false
	carJSON, err := json.Marshal(carDy)
	if err != nil {
		return "", err
	}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "unlock", Args: [][]byte{carJSON}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		log.Printf("request err:%s !\n", err)
		return "", err
	}

	return "successful!", nil
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

func UpdataRoad(roadInfo *def.RoadInformation) (string, error) {
	rInfoJSON, err := json.Marshal(roadInfo)
	if err != nil {
		return "", err
	}

	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "updataRoad", Args: [][]byte{rInfoJSON, []byte(eventID)}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		return "", err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return resp, nil
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

func CheckCollision(onroad *def.OnRoad) ([]byte, error) {
	onroadASJSON, err := json.Marshal(onroad)
	if err != nil {
		return []byte{}, err
	}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "checkCollision", Args: [][]byte{onroadASJSON}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}

	return []byte("seccessful!"), nil
}

func CheckRGL(onroad *def.OnRoad) ([]byte, error) {
	onroadASJSON, err := json.Marshal(onroad)
	if err != nil {
		return []byte{}, err
	}

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "checkRGL", Args: [][]byte{onroadASJSON}}
	_, err = ServiceClient.Client.Execute(req)
	if err != nil {
		return []byte{}, err
	}

	return []byte("seccessful!"), nil
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

func DealRGL(rglId string) (string, error) {
	eventID, _ := utils.NewUUID()
	reg, notifier := regitserEvent(ServiceClient.Client, ServiceClient.ChaincodeID, eventID)
	defer ServiceClient.Client.UnregisterChaincodeEvent(reg)

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "deleteRGL", Args: [][]byte{[]byte(rglId), []byte(eventID)}}
	_, err := ServiceClient.Client.Execute(req)
	if err != nil {
		return "", err
	}

	resp, err := eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func GetHistoryRGL(rglId string) ([]byte, error) {

	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "getHistoryRGL", Args: [][]byte{[]byte(rglId)}}
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}

	return resp.Payload, nil
}

func CarRGL(carNum string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "carRGL", Args: [][]byte{[]byte(carNum)}}
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}

	return resp.Payload, nil
}

func getLock(carNum string) ([]byte, error) {
	req := channel.Request{ChaincodeID: ServiceClient.ChaincodeID, Fcn: "getlock", Args: [][]byte{[]byte(carNum)}}
	resp, err := ServiceClient.Client.Query(req)
	if err != nil {
		return []byte{}, err
	}

	return resp.Payload, nil
}

// //锁监听
// func LockCheck(carNum string, conn *net.TCPConn) {
// 	resp, err := getLock(carNum)
// 	carDy := &def.Lock{}

// 	if err = json.Unmarshal(resp, carDy); err != nil {
// 		log.Printf("监听 unmarshal errror :%s\n", err)
// 	}

// 	if carDy.Lc {
// 		sendTips(conn, "car was locked!\n")
// 	} else {
// 		sendTips(conn, "car was opened!\n")
// 	}
// }

//错误码监听
func FaulCheck(carNum string, conn *net.TCPConn) {
	resp, err := GetState(carNum)
	carDy := &def.CarDy{}

	if err = json.Unmarshal(resp, carDy); err != nil {
		log.Printf("监听 unmarshal errror :%s\n", err)
	}

	if carDy.FaultCode != "" {
		sendTips(conn, "has a fault code :"+carDy.FaultCode+"!\n")
	}
}

// //违章监听
// func RglCheck(carNum string, conn *net.TCPConn) {
// 	resp, err := CarRGL(carNum)
// 	rglInfoItem := &def.CarRGLItem{}

// 	if err = json.Unmarshal(resp, rglInfoItem); err != nil {
// 		log.Printf("监听 unmarshal errror :%s\n", err)
// 	}

// 	if rglInfoItem != nil {
// 		sendTips(conn, "has a Regulations info :"+rglInfoItem.Item[0].Mes+"!\n")
// 	}
// }

// func CollisCheck(carNum string, conn *net.TCPConn) {
// 	sendTips(conn, "有追尾风险！")
// }

//onRoad listen
// func OnRoadListen(carNum string, conn *net.TCPConn) {
// 	// CheckRGL(cGRL *def.CheckRGL)
// 	sendTips(conn, "on the tag 1 road!")
// }

func sendTips(conn *net.TCPConn, message string) {
	log.Println("send tips:" + message)
	mes := &def.SocketResp{Mes: message}

	resp, err := json.Marshal(mes)
	if err != nil {
		log.Println("send tip err:" + err.Error())
		return
	}
	conn.Write(resp)
	conn.Write([]byte("\n"))
}

var MessageMap *sync.Map

func init() {
	MessageMap = &sync.Map{}
}

//锁监听
func LockCheck(carNum string) {
	resp, err := getLock(carNum)
	carDy := &def.Lock{}

	if err = json.Unmarshal(resp, carDy); err != nil {
		log.Printf("监听 unmarshal errror :%s\n", err)
	}

	if carDy.Lc {
		MessageMap.Store("lock", "car was locked!")
	} else {
		MessageMap.Store("lock", "car was opened!")
	}
}

//违章监听
func RglCheck(carNum string) {
	resp, err := CarRGL(carNum)
	rglInfoItem := &def.CarRGLItem{}

	if err = json.Unmarshal(resp, rglInfoItem); err != nil {
		log.Printf("监听 unmarshal errror :%s\n", err)
	}

	if rglInfoItem != nil {
		MessageMap.Store("rgl", "has a Regulations info :"+rglInfoItem.Item[0].Mes+"!\n")
	}
}

func CollisCheck(carNum string) {
	MessageMap.Store("col", "有追尾风险！")
}
