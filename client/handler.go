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
		return
	}

	resp, err := service.InitCar(ubody)
	if err != nil {
		fmt.Printf("servic init error 2:%s\n", err)
	}

	fmt.Printf("Car init successful:%s!\n", resp)
}

//request body:string car id
func GetCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	carNum := p.ByName("car_num")

	car, err := service.GetCar(carNum)
	if err != nil {
		log.Printf("error:%s\n", err)
		return
	}

	fmt.Println(car)
}

//request body:json carDyReq
func PutCarDy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.CarDyReq{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println("request error!")
		return
	}

	err := service.PutCarDy(ubody)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("put carDy successful!\n")
}

//request body: json carDyReq
func LockCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.CarDyReq{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println("request error!")
		return
	}

	resp, err := service.LockCar(ubody)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("lock car successful:%s!\n", resp)
}

//request : owner
func QueryCarByOwner(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	o := p.ByName("owner_name")
	resp, err := service.QueryCarByOwner(o)
	if err != nil {
		log.Printf("error:%s\n", err)
		return
	}

	fmt.Println(resp)
}

//request: carNum
func QueryCarHistry(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	carNum := p.ByName("car_num")
	resp, err := service.QueryHistoryForCar(carNum)
	if err != nil {
		log.Printf("error:%s\n", err)
		return
	}

	fmt.Println(resp)
}
