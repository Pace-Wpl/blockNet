package main

import (
	"log"
	"net"
	"net/http"
	"time"

	ds "github.com/blockNet/client/daemontask"
	"github.com/blockNet/client/service"
	"github.com/julienschmidt/httprouter"
)

type tcpServer struct {
	Chan chan string
	Conn *net.TCPConn
}

var tcps tcpServer

type middleWare struct {
	r *httprouter.Router
}

func NewMiddleWare(r *httprouter.Router) http.Handler {
	m := middleWare{}
	m.r = r
	return m
}

func (m middleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("request ip:%s\n", r.RemoteAddr)
	// if strings.Split(r.RemoteAddr, ":")[0] != "127.0.0.1" {
	// 	w.Write([]byte("bad request."))
	// 	return
	// }
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/init", InitCar)
	router.GET("/car/:car_id", GetCar)
	router.POST("/carDy", PutCarDy)
	router.POST("/lcar", LockCar)
	router.POST("/ulcar", UnLockCar)
	router.GET("/carDyh/:car_num", GetCarHistory)
	router.POST("/onRoad", OnRoad)
	router.POST("/rgl", GetRGL)
	router.GET("/carRgl/:car_num", CarRGL)
	router.POST("/deaRgl", DealRGL)
	router.POST("/rglH", GetRglHistory)
	router.POST("/road", PutRoad)
	router.GET("/task/:car_num", TaskOpen)

	return router
}

func main() {
	tcps.Chan = make(chan string)
	go tcps.tcpListen()
	defer func() {
		service.ServiceClient.Sdk.Close()
		tcps.Chan <- "done"
		time.Sleep(time.Duration(10) * time.Second)
	}()

	r := RegisterHandlers()
	m := NewMiddleWare(r)
	http.ListenAndServe(":8000", m)

}

func (t *tcpServer) tcpListen() {
	var tcpAddr *net.TCPAddr
	address := "127.0.0.1:9999"
	tcpAddr, _ = net.ResolveTCPAddr("tcp", address)
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer func() {
		log.Println("tcp close!")
		tcpListener.Close()
	}()
	log.Println("Server ready to read ...")

	tcpConn, err := tcpListener.AcceptTCP()
	if err != nil {
		log.Println(err)
		return
	}

	t.Conn = tcpConn
	log.Println("A client connected :" + tcpConn.RemoteAddr().String())

forloop:
	for {
		select {
		case flag := <-t.Chan:
			if flag == "done" {
				break forloop
			} else {
				task := ds.NewTask(tcps.Conn)
				task.StartDaemon(flag)
				defer task.StopDaemon()
			}

		}
	}

	select {}

}
