package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/blockNet/clientTcp/def"
)

//消息处理
func register(conn *net.TCPConn, msg string) {
	com := strings.Split(msg, ";")
	switch com[0] {
	case "getCar":
		GetCar(conn, com[1:])
	}

}

func main() {
	//定义一个tcp断点
	var tcpAddr *net.TCPAddr
	address := def.SERVER_IP + ":" + strconv.Itoa(def.SERVER_PORT)
	tcpAddr, _ = net.ResolveTCPAddr("tcp", address)
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()
	fmt.Println("Server ready to read ...")
	//循环接收客户端的连接，创建一个协程具体去处理连接
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("A client connected :" + tcpConn.RemoteAddr().String())
		if strings.Split(tcpConn.RemoteAddr().String(), ":")[0] != "127.0.0.1" {
			fmt.Println([]byte("bad request."))
			continue
		} else {
			fmt.Println([]byte("connecte and server!"))
			go tcpHandle(tcpConn)
		}

	}
}

//处理连接
func tcpHandle(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()

	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil || err == io.EOF || message == "END" {
			break
		}
		register(conn, string(message))
	}
}
