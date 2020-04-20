package main

import (
	"net"

	"github.com/blockNet/clientTcp/def"
)

func sendErrorResponse(conn *net.TCPConn, e error) {
	conn.Write([]byte(e.Error()))
}

func sendNormalResponse(conn *net.TCPConn, resp *def.Response) {
	conn.Write(resp.Msg)
}
