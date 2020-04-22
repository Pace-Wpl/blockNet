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

func InitCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.CarInit{}

	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := service.InitCar(ubody); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Car init successful!\n")
}

func GetCar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	carNum := p.ByName("car_num")

	car, err := service.GetCar(carNum)
	if err != nil {
		log.Printf("error:%s\n", err)
		return
	}

	fmt.Println(car)
}
