package main

import (
	"bytes"
	"log"
	"net"
)

func sendErrorResponse(conn *net.TCPConn, e error) {
	conn.Write([]byte(e.Error() + "\n"))
}

func sendNormalResponse(conn *net.TCPConn, mes []byte) {
	var b bytes.Buffer
	b.Write(mes)
	b.Write([]byte("\n"))
	log.Println(string(b.Bytes()))
	conn.Write(b.Bytes())
}
