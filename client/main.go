package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/car", InitCar)
	router.GET("/car/:car_num", GetCar)

	return router
}

func main() {

	r := RegisterHandlers()
	http.ListenAndServe(":8000", r)
}
