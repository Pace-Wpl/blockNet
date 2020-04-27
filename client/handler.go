package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blockNet/client/def"
	"github.com/blockNet/client/service"
	"github.com/julienschmidt/httprouter"
)

func OnRoad(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	log.Println(string(res))
	ubody := &def.OnRoadReq{}

	if err := json.Unmarshal(res, ubody); err != nil {
		log.Println(err)
		sendErrorResponse(w, def.ERROR_BAD_REQUETS)
		return
	}

	onRoad := &def.OnRoad{
		ObjectType: "onRoad", CarNum: ubody.CarNum, Velocity: ubody.Velocity, Direction: ubody.Direction,
		Position: ubody.Position, Code: ubody.Code,
	}

	_, err := service.CheckCollision(onRoad)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	resp, err := service.CheckRGL(onRoad)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	sendNormalResponse(w, string(resp))
}

func GetRGL(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.RglReq{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, def.ERROR_BAD_REQUETS)
		return
	}
	resp, err := service.GetRGL(ubody.RglID)
	if err != nil {
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	sendNormalResponse(w, string(resp))
}

func CarRGL(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	carNum := p.ByName("car_num")

	resp, err := service.CarRGL(carNum)
	if err != nil {
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	sendNormalResponse(w, string(resp))

}

func DealRGL(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.RglReq{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, def.ERROR_BAD_REQUETS)
		return
	}
	resp, err := service.DealRGL(ubody.RglID)
	if err != nil {
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	sendNormalResponse(w, resp)
}

func GetRglHistory(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.RglReq{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, def.ERROR_BAD_REQUETS)
		return
	}
	resp, err := service.GetHistoryRGL(ubody.RglID)
	if err != nil {
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	sendNormalResponse(w, string(resp))
}

func PutRoad(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.RoadInformation{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, def.ERROR_BAD_REQUETS)
		return
	}
	resp, err := service.UpdataRoad(ubody)
	if err != nil {
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	sendNormalResponse(w, resp)

}

func TestChaincode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k := p.ByName("key")
	fmt.Println(k)
	if err := service.TestChaincod(k); err != nil {
		fmt.Println(err.Error())
		return
	}

}

func TestChaincodQ(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k := p.ByName("key")
	fmt.Println(k)
	if err := service.TestChaincodQ(k); err != nil {
		fmt.Println(err.Error())
		return
	}
}

//initCar;
//request body:json carinit
func InitCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.CarInit{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Printf("unmarshal error 1:%s\n", err)
		sendErrorResponse(w, def.ERROR_BAD_REQUETS)
		return
	}

	resp, err := service.InitCar(ubody)
	if err != nil {
		fmt.Printf("servic init error 2:%s\n", err)
		sendErrorResponse(w, def.ERROR_INTERNAL)
	}

	tcps.Chan <- ubody.CarNumber //开启TCP任务

	sendNormalResponse(w, resp)
}

//request body:string car id
func GetCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	carNum := p.ByName("car_id")

	car, err := service.GetCar(carNum)
	if err != nil {
		log.Printf("error:%s\n", err)
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	sendNormalResponse(w, string(car))
}

//request body:json carDyReq
func PutCarDy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	log.Println(string(res))
	ubody := &def.CarDyReq{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, def.ERROR_BAD_REQUETS)
		return
	}

	err := service.PutCarDy(ubody)
	if err != nil {
		fmt.Println(err)
	}

	sendNormalResponse(w, "successful!")
}

//request body: json carDyReq
func LockCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.LockCar{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println("request error!")
		return
	}

	resp, err := service.LockCar(ubody)
	if err != nil {
		fmt.Println(err)
	}

	sendNormalResponse(w, resp)
}

func UnLockCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.LockCar{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println("request error!")
		return
	}

	resp, err := service.UnLockCar(ubody)
	if err != nil {
		fmt.Println(err)
	}

	sendNormalResponse(w, resp)
}

//request : owner
// func QueryCarByOwner(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	o := p.ByName("owner_name")
// 	resp, err := service.QueryCarByOwner(o)
// 	if err != nil {
// 		log.Printf("error:%s\n", err)
// 		return
// 	}

// 	fmt.Println(resp)
// }

//request: carNum
func GetCarHistory(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	carNum := p.ByName("car_num")
	resp, err := service.QueryHistoryForCar(carNum)
	if err != nil {
		log.Printf("error:%s\n", err)
		return
	}

	sendNormalResponse(w, string(resp))
}

func TaskOpen(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	carNum := p.ByName("car_num")
	tcps.Chan <- carNum //开启TCP任务
	sendNormalResponse(w, "open task:"+carNum)
}
