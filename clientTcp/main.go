package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	ds "github.com/blockNet/clientTcp/daemontask"
	"github.com/blockNet/clientTcp/def"
)

//消息处理
func tcpHandle(conn *net.TCPConn, com []string) {
	switch com[0] {
	case "initCar":
		initCar(conn, com[1:])
	case "getCar":
		getCar(conn, com[1:])
	case "getState":
		getState(conn, com[1:])
	case "putCarDy":
		putCarDy(conn, com[1:])
	case "lockCar":
		lockCar(conn, com[1:])
	case "unlock":
		unLockCar(conn, com[1:])
	case "queryCarByOwner":
		queryCarByOwner(conn, com[1:])
	case "queryCarHistry":
		queryCarHistry(conn, com[1:])
	case "deleteCar":
		deleteCar(conn, com[1:])
	case "checkRGL":
		checkRGL(conn, com[1:])
	case "updataCar":
		updataCar(conn, com[1:])
	case "updataRoad":
		updataRoad(conn, com[1:])
	case "readRoad":
		getRoad(conn, com[1:])
	case "deleteRoad":
		deleteRoad(conn, com[1:])
	case "onRoad":
		onRoad(conn, com[1:])
	case "updataRGL":
		updataRGL(conn, com[1:])
	case "readRGL":
		getRGL(conn, com[1:])
	case "carRGL":
		carRGL(conn, com[1:])
	case "deleteRGL":
		dealRGL(conn, com[1:])
	case "getHistoryRGL":
		getHistoryRGL(conn, com[1:])
	default:
		sendErrorResponse(conn, errors.New("没有相应的方法！"))
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
			fmt.Println("bad request.")
			continue
		} else {
			fmt.Println("connecte and server!")
			go register(tcpConn)
		}

	}
}

//处理连接
func register(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()

	conn.Write([]byte("connect successful!\n"))
	c := make(chan string)
	t := &ds.Task{Controller: c, Conn: conn}

	defer func() {
		t.StopDaemon()
		fmt.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		com := strings.Split(message, ",")
		fmt.Println("com0:" + com[0])
		if err != nil || err == io.EOF || com[0] == "EXIT" {
			fmt.Println("exit and send exit")
			conn.Write([]byte("EXIT"))
			break
		} else if com[0] == "TASK" {
			task := com[1]
			if err != nil || err == io.EOF || task == "EXIT" {
				break
			}
			t.StartDaemon(task)
		}
		tcpHandle(conn, com)
	}
}
