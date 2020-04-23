package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/blockNet/clientTcp/def"
	"github.com/blockNet/clientTcp/service"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
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
//request args: name,carnum,id,owner,type,colour,lock,commander,velociy,temperature,faultcode
func initCar(conn *net.TCPConn, msgs []string) {
	lock, err := strconv.ParseBool(msgs[6])
	if err != nil {
		sendErrorResponse(conn, errors.New("参数lock,8错误！"))
		return
	}

	v, err := strconv.ParseFloat(msgs[8], 32)
	if err != nil {
		sendErrorResponse(conn, errors.New("参数velocity,10错误！"))
		return
	}

	t, err := strconv.ParseFloat(msgs[9], 32)
	if err != nil {
		sendErrorResponse(conn, errors.New("参数Temperature,11错误！"))
		return
	}

	ubody := &def.CarInit{
		Name: msgs[0], CarNumber: msgs[1], ID: msgs[2], Owner: msgs[3], Type: msgs[4], Colour: msgs[5],
		Lock: lock, Commander: msgs[7], Velocity: float32(v), Temperature: float32(t), FaultCode: msgs[10],
	}

	resp, err := service.InitCar(ubody)
	if err != nil {
		fmt.Printf("servic init error 2:%s\n", err)
		sendErrorResponse(conn, err)
	}

	sendNormalResponse(conn, []byte(resp))
}

//request args:car id
func GetCar(conn *net.TCPConn, args []string) {
	carNum := args[0]

	car, err := service.GetCar(carNum)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, car)
}

//request args:carnum
func GetState(conn *net.TCPConn, args []string) {
	carNum := args[0]

	car, err := service.GetState(carNum)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, car)
}

//request body:carnum,lock,commander,v,t,f
func PutCarDy(conn *net.TCPConn, args []string) {
	lock, err := strconv.ParseBool(args[1])
	if err != nil {
		sendErrorResponse(conn, errors.New("参数lock,3错误！"))
		return
	}

	v, err := strconv.ParseFloat(args[3], 32)
	if err != nil {
		sendErrorResponse(conn, errors.New("参数velocity,5错误！"))
		return
	}

	t, err := strconv.ParseFloat(args[4], 32)
	if err != nil {
		sendErrorResponse(conn, errors.New("参数Temperature,6错误！"))
		return
	}
	ubody := &def.CarDyReq{
		CarNumber: args[0], Lock: lock, Commander: args[2], Velocity: float32(v), Temperature: float32(t),
		FaultCode: args[5],
	}

	err = service.PutCarDy(ubody)
	if err != nil {
		sendErrorResponse(conn, err)
		fmt.Println(err)
	}

	sendNormalResponse(conn, []byte("Successful !"))
}

//request body: carnum,lock,commander,v,t,f
func LockCar(conn *net.TCPConn, args []string) {
	lock, err := strconv.ParseBool(args[1])
	if err != nil {
		sendErrorResponse(conn, errors.New("参数lock,3错误！"))
		return
	}

	v, err := strconv.ParseFloat(args[3], 32)
	if err != nil {
		sendErrorResponse(conn, errors.New("参数velocity,5错误！"))
		return
	}

	t, err := strconv.ParseFloat(args[4], 32)
	if err != nil {
		sendErrorResponse(conn, errors.New("参数Temperature,6错误！"))
		return
	}
	ubody := &def.CarDyReq{
		CarNumber: args[0], Lock: lock, Commander: args[2], Velocity: float32(v), Temperature: float32(t),
		FaultCode: args[5],
	}

	resp, err := service.LockCar(ubody)
	if err != nil {
		sendErrorResponse(conn, err)
		fmt.Println(err)
	}

	sendNormalResponse(conn, []byte(resp))
}

//request : owner
func QueryCarByOwner(conn *net.TCPConn, args []string) {
	o := args[0]
	resp, err := service.QueryCarByOwner(o)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//request: carNum
func QueryCarHistry(conn *net.TCPConn, args []string) {
	carNum := args[0]
	resp, err := service.QueryHistoryForCar(carNum)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}
