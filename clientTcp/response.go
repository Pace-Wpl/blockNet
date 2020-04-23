package main

import (
	"net"
)

func sendErrorResponse(conn *net.TCPConn, e error) {
	conn.Write([]byte(e.Error()))
}

func sendNormalResponse(conn *net.TCPConn, mes []byte) {
	conn.Write(mes)
}
