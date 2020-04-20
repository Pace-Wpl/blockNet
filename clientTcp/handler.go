package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/blockNet/clientTcp/def"
	"github.com/blockNet/clientTcp/service"
)

func GetCar(conn *net.TCPConn, parm []string) {
	car, err := service.GetCar(parm[0])
	if err != nil {
		log.Printf("error:%s\n", err)
		sendErrorResponse(conn, err)
		return
	}

	c, err := json.Marshal(car)
	if err != nil {
		log.Printf("error:%s\n", err)
		sendErrorResponse(conn, err)
		return
	}

	resp := &def.Response{Msg: c}
	sendNormalResponse(conn, resp)
}
