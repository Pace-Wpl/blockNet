package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/blockNet/clientTcp/def"
)

//消息处理
func tcpHandle(conn *net.TCPConn, msg string) {
	com := strings.Split(msg, ",")
	switch com[0] {
	case "initCar":
		initCar(conn, com[1:])
	case "getCar":
		GetCar(conn, com[1:])
	case "getState":
		GetState(conn, com[1:])
	case "putCarDy":
		PutCarDy(conn, com[1:])
	case "lockCar":
		LockCar(conn, com[1:])
	case "unlock":
		return t.unLockCar(stub, args)
	case "queryCarByOwner":
		QueryCarByOwner(conn, com[1:])
	case "queryCarHistry":
		QueryCarHistry(conn, com[1:])
	case "deleteCar":
		return t.deleteCar(stub, args)
	case "checkRGL":
		return t.checkRGL(stub, args)
	case "updataCar":
		return t.updataCar(stub, args)
	case "updataRoad":
		return t.updataRoad(stub, args)
	case "readRoad":
		return t.readRoad(stub, args)
	case "deleteRoad":
		return t.deleteRoad(stub, args)
	case "onRoad":
		return t.onRoad(stub, args)
	case "updataRGL":
		return t.updataRGL(stub, args)
	case "readRGL":
		return t.readRgl(stub, args)
	case "deleteRGL":
		return t.dealRGL(stub, args)
	case "getHistoryRGL":
		return t.getHistryRgl(stub, args)
	default:
	}
	sendErrorResponse(conn, errors.New("没有相应的方法！"))
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
			go register(tcpConn)
		}

	}
}

//处理连接
func register(conn *net.TCPConn) {
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
		tcpHandle(conn, string(message))
	}
}
