package main

import (
	"log"
	"net/http"
	"strings"

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
	log.Printf("request ip:%s\n", r.RemoteAddr)
	if strings.Split(r.RemoteAddr, ":")[0] != "127.0.0.1" {
		w.Write([]byte("bad request."))
		return
	} else {
		w.Write([]byte("auth pass!"))
	}
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/car", InitCar)
	router.GET("/car/:car_num", GetCar)
	router.GET("/test/:key", TestChaincode)
	router.GET("/q/:key", TestChaincodQ)

	return router
}

func main() {
	defer service.ServiceClient.Sdk.Close()
	r := RegisterHandlers()
	m := NewMiddleWare(r)
	http.ListenAndServe(":8000", m)

}
