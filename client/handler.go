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
	rglID := p.ByName("rgl_id")

	resp, err := service.GetRGL(rglID)
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
	rglID := p.ByName("rgl_id")

	resp, err := service.DealRGL(rglID)
	if err != nil {
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	sendNormalResponse(w, resp)
}

func GetRglHistory(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rglID := p.ByName("rgl_id")

	resp, err := service.GetHistoryRGL(rglID)
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
	log.Println(string(res))
	ubody := &def.CarInit{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Printf("unmarshal error 1:%s\n", err)
		sendErrorResponse(w, def.ERROR_BAD_REQUETS)
		return
	}

	resp, err := service.InitCar(ubody)
	if err != nil {
		fmt.Printf("servic init error 2:%s\n", err)
		sendBadResponse(w, "user 不存在！")
		return
	}

	go task.StartDaemon(ubody.CarNumber)

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
		return
	}

	sendNormalResponse(w, "successful!")
}

//request body: json carDyReq
func LockCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.LockCarReq{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println("request error!")
		return
	}

	lock := &def.LockCar{ObjectType: "carLock", CarNum: ubody.CarNum, Lc: true, Certificate: ubody.Certificate}

	resp, err := service.LockCar(lock)
	if err != nil {
		fmt.Println(err)
		sendBadResponse(w, "Unauthorized")
		return
	}

	sendNormalResponse(w, resp)
}

func UnLockCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.LockCarReq{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println("request error!")
		return
	}

	lock := &def.LockCar{ObjectType: "carLock", CarNum: ubody.CarNum, Lc: false, Certificate: ubody.Certificate}
	resp, err := service.UnLockCar(lock)
	if err != nil {
		sendBadResponse(w, "Unauthorized")
		return
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
		sendErrorResponse(w, def.ERROR_INTERNAL)
		return
	}

	sendNormalResponse(w, string(resp))
}

//equest body: json UserReg
func Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.UserReg{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println("request error!")
		sendErrorResponse(w, "")
		return
	}

	cer, err := service.Register(ubody)
	if err != nil {
		fmt.Println(err)
		sendBadResponse(w, "username exit!")
		return
	}

	sendNormalResponse(w, "注册成功:"+cer)
}

//equest body: json UserReg
func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.UserReg{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println("request error!")
		sendErrorResponse(w, def.ERROR_BAD_REQUETS)
		return
	}

	_, err := service.GetUser(ubody)
	if err != nil {
		fmt.Println(err)
		sendBadResponse(w, "password error!")
		return
	}

	sendNormalResponse(w, "欢迎您:"+ubody.UserName)
}

func TaskOpen(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	carNum := p.ByName("car_num")
	go task.StartDaemon(carNum)
	sendNormalResponse(w, "open task:"+carNum)
}

func Lms(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	mes, ok := service.MessageMap.Load("lock")
	if ok {
		log.Println("send mes :" + mes.(string))
		sendNormalResponse(w, mes.(string))
		service.MessageMap.Delete("lock")
	}
	mes1, ok1 := service.MessageMap.Load("col")
	if ok1 {
		log.Println("send mes :" + mes1.(string))
		sendNormalResponse(w, mes1.(string))
		service.MessageMap.Delete("col")
	}
	mes2, ok2 := service.MessageMap.Load("rgl")
	if ok2 {
		log.Println("send mes :" + mes2.(string))
		sendNormalResponse(w, mes2.(string))
		service.MessageMap.Delete("rgl")
	}
	// sendNormalResponse(w, "")
}
