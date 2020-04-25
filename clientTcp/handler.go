package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/blockNet/clientTcp/def"
	"github.com/blockNet/clientTcp/service"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func testChaincode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k := p.ByName("key")
	fmt.Println(k)
	if err := service.TestChaincod(k); err != nil {
		fmt.Println(err.Error())
		return
	}

}

func testChaincodQ(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
func getCar(conn *net.TCPConn, args []string) {
	carNum := args[0]
	fmt.Println(carNum)
	car, err := service.GetCar(carNum)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, car)
}

//request args:carnum
func getState(conn *net.TCPConn, args []string) {
	carNum := args[0]

	car, err := service.GetState(carNum)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, car)
}

//request body:carnum,lock,commander,v,t,f
func putCarDy(conn *net.TCPConn, args []string) {
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

//request body: carnum,lock,commander,(v,t,f)
func lockCar(conn *net.TCPConn, args []string) {
	lock, err := strconv.ParseBool(args[1])
	if err != nil {
		sendErrorResponse(conn, errors.New("参数lock,3错误！"))
		return
	}

	// v, err := strconv.ParseFloat(args[3], 32)
	// if err != nil {
	// 	sendErrorResponse(conn, errors.New("参数velocity,5错误！"))
	// 	return
	// }

	// t, err := strconv.ParseFloat(args[4], 32)
	// if err != nil {
	// 	sendErrorResponse(conn, errors.New("参数Temperature,6错误！"))
	// 	return
	// }
	ubody := &def.CarDyReq{
		CarNumber: args[0], Lock: lock, Commander: args[2], Velocity: 0, Temperature: 28,
		FaultCode: "",
	}

	resp, err := service.LockCar(ubody)
	if err != nil {
		sendErrorResponse(conn, err)
		fmt.Println(err)
	}

	sendNormalResponse(conn, []byte(resp))
}

//request body: carnum,lock,commander,(v,t,f)
func unLockCar(conn *net.TCPConn, args []string) {
	lock, err := strconv.ParseBool(args[1])
	if err != nil {
		sendErrorResponse(conn, errors.New("参数lock,3错误！"))
		return
	}

	// v, err := strconv.ParseFloat(args[3], 32)
	// if err != nil {
	// 	sendErrorResponse(conn, errors.New("参数velocity,5错误！"))
	// 	return
	// }

	// t, err := strconv.ParseFloat(args[4], 32)
	// if err != nil {
	// 	sendErrorResponse(conn, errors.New("参数Temperature,6错误！"))
	// 	return
	// }
	ubody := &def.CarDyReq{
		CarNumber: args[0], Lock: lock, Commander: args[2], Velocity: 0, Temperature: 28,
		FaultCode: "",
	}

	resp, err := service.UnLockCar(ubody)
	if err != nil {
		sendErrorResponse(conn, err)
		fmt.Println(err)
	}

	sendNormalResponse(conn, []byte(resp))
}

//request : owner
func queryCarByOwner(conn *net.TCPConn, args []string) {
	o := args[0]
	resp, err := service.QueryCarByOwner(o)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	carItem := &def.OwenrCarItem{}
	if err = json.Unmarshal(resp, carItem); err != nil {
		sendErrorResponse(conn, err)
		return
	}

	var buffer bytes.Buffer
	isWrite := false
	for _, v := range carItem.Item {
		v1, _ := json.Marshal(v)
		if !isWrite {
			buffer.WriteString(v.Owner + ":\n")
			isWrite = true
		}
		buffer.Write(v1)
		buffer.WriteString(";\n")
	}

	sendNormalResponse(conn, buffer.Bytes())
}

//request: carNum
func queryCarHistry(conn *net.TCPConn, args []string) {
	carNum := args[0]
	resp, err := service.QueryHistoryForCar(carNum)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	hisItem := &def.HistryItem{}
	if err = json.Unmarshal(resp, hisItem); err != nil {
		sendErrorResponse(conn, err)
		return
	}

	var buffer bytes.Buffer
	isWrite := false
	for _, v := range hisItem.Item {
		v1, _ := json.Marshal(v.CarDy)
		if !isWrite {
			buffer.WriteString(carNum + ":\n")
			isWrite = true
		}
		buffer.WriteString("Timestamp:" + v.Timestamp + ";\n")
		buffer.WriteString("IsDelete:" + strconv.FormatBool(v.IsDelete) + ";\n")
		buffer.Write(v1)
		buffer.WriteString(";\n")
	}

	sendNormalResponse(conn, buffer.Bytes())
}

//reuqets: carid
func deleteCar(conn *net.TCPConn, args []string) {
	carID := args[0]
	resp, err := service.DeleteCar(carID)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//reuqest: carnum,road,v
func checkRGL(conn *net.TCPConn, args []string) {
	v, err := strconv.ParseFloat(args[2], 32)
	if err != nil {
		sendErrorResponse(conn, errors.New("参数4错误！"))
		return
	}
	cRgl := &def.CheckRGL{
		CarNumber: args[0], Road: args[1], Velocity: float32(v),
	}

	resp, err := service.CheckRGL(cRgl)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//request: name,carnum,carid,owner,type,colour
func updataCar(conn *net.TCPConn, args []string) {
	carInfo := &def.CarInfomation{
		ObjectType: "carInfomation", Name: args[0], CarNumber: args[1], ID: args[2], Owner: args[3],
		Type: args[4], Colour: args[5],
	}

	resp, err := service.UpdataCar(carInfo)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

///test above
//request: road code,name,Coordinate,type,limit,tag
func updataRoad(conn *net.TCPConn, args []string) {
	l, err := strconv.ParseFloat(args[4], 32)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	t, err := strconv.Atoi(args[5])
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	roadInfo := &def.RoadInformation{
		ObjectType: "roadInfomation", Code: args[0], Name: args[1], Coordinate: args[2], Type: args[3],
		Limit: float32(l), Tag: t,
	}

	resp, err := service.UpdataRoad(roadInfo)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//request: road code
func getRoad(conn *net.TCPConn, args []string) {
	roadCode := args[0]

	resp, err := service.GetRoad(roadCode)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//request: road code
func deleteRoad(conn *net.TCPConn, args []string) {
	roadCode := args[0]

	resp, err := service.GetRoad(roadCode)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//request: road code,carnum,v
func onRoad(conn *net.TCPConn, args []string) {
	roadCode := args[0]
	carNum := args[1]
	v, err := strconv.ParseFloat(args[2], 32)
	if err != nil {
		sendErrorResponse(conn, err)
	}

	resp, err := service.OnRoad(roadCode, carNum)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	cGRL := &def.CheckRGL{CarNumber: carNum, Road: roadCode, Velocity: float32(v)}
	_, err = service.CheckRGL(cGRL)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//request:rgl id,carnum,road,type,mes
func updataRGL(conn *net.TCPConn, args []string) {
	rglInfo := &def.RegulationsInfo{
		ObjectType: "RegulationsInfo", ID: args[0], CarNumber: args[1], Road: args[2], Type: args[3],
		Mes: args[4],
	}

	resp, err := service.UpdataRGL(rglInfo)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//request:rgl id
func getRGL(conn *net.TCPConn, args []string) {
	rglId := args[0]

	resp, err := service.GetRGL(rglId)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//request:rgl id
func dealRGL(conn *net.TCPConn, args []string) {
	rglId := args[0]

	resp, err := service.DealRGL(rglId)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	sendNormalResponse(conn, resp)
}

//request:rgl id
func getHistoryRGL(conn *net.TCPConn, args []string) {
	rglId := args[0]

	resp, err := service.GetHistoryRGL(rglId)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	hisItem := &def.HistryRGLItem{}
	if err = json.Unmarshal(resp, hisItem); err != nil {
		sendErrorResponse(conn, err)
		return
	}

	var buffer bytes.Buffer
	isWrite := false
	for _, v := range hisItem.Item {
		v1, _ := json.Marshal(v.Rgl)
		if !isWrite {
			buffer.WriteString(rglId + ":\n")
			isWrite = true
		}
		buffer.WriteString("Timestamp:" + v.Timestamp + ";\n")
		buffer.WriteString("IsDelete:" + strconv.FormatBool(v.IsDelete) + ";\n")
		buffer.Write(v1)
		buffer.WriteString(";\n")
	}

	sendNormalResponse(conn, buffer.Bytes())
}

//request: carnum
func carRGL(conn *net.TCPConn, args []string) {
	carNum := args[0]

	resp, err := service.CarRGL(carNum)
	if err != nil {
		sendErrorResponse(conn, err)
		return
	}

	rglItem := &def.CarRGLItem{}
	if err = json.Unmarshal(resp, rglItem); err != nil {
		sendErrorResponse(conn, err)
		return
	}

	var buffer bytes.Buffer
	isWrite := false
	for _, v := range rglItem.Item {
		v1, _ := json.Marshal(v)
		if !isWrite {
			buffer.WriteString(carNum + ":\n")
			isWrite = true
		}
		buffer.Write(v1)
		buffer.WriteString(";\n")
	}

	sendNormalResponse(conn, buffer.Bytes())
}
