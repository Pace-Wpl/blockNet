package main

import (
	"net/http"
	"time"

	ds "github.com/blockNet/client/daemontask"
	"github.com/blockNet/client/service"
	"github.com/julienschmidt/httprouter"
)

type middleWare struct {
	r *httprouter.Router
}

func NewMiddleWare(r *httprouter.Router) http.Handler {
	m := middleWare{}
	m.r = r
	return m
}

func (m middleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log.Printf("request ip:%s\n", r.RemoteAddr)
	// if strings.Split(r.RemoteAddr, ":")[0] != "127.0.0.1" {
	// 	w.Write([]byte("bad request."))
	// 	return
	// }
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type,Depth, User-Agent, X-Session-Id, X-User-Name, If-Modified-Since, Cache-Control, Origin")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("access-control-allow-methods", "GET, POST, OPTIONS, PUT, DELETE")
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/car", InitCar)
	router.GET("/car/:car_id", GetCar)
	router.POST("/carDy", PutCarDy)
	router.POST("/lcar", LockCar)
	router.POST("/ulcar", UnLockCar)
	router.GET("/carDyh/:car_num", GetCarHistory)
	router.POST("/onRoad", OnRoad)
	router.GET("/rgl/:rgl_id", GetRGL)
	router.GET("/carRgl/:car_num", CarRGL)
	router.GET("/deaRgl/:rgl_id", DealRGL)
	router.GET("/rglH/:rgl_id", GetRglHistory)
	router.POST("/road", PutRoad)
	router.GET("/task/:car_num", TaskOpen)
	router.GET("/lms", Lms)

	return router
}

var task ds.Task

func main() {
	Ch := make(chan string)
	task.Controller = Ch
	defer func() {
		service.ServiceClient.Sdk.Close()
		task.StopDaemon()
		time.Sleep(time.Duration(10) * time.Second)
	}()

	r := RegisterHandlers()
	m := NewMiddleWare(r)
	http.ListenAndServe(":8000", m)

}
