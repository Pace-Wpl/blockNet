package main

import (
	"io"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, resp string) {
	w.WriteHeader(400)
	io.WriteString(w, resp)
}

func sendNormalResponse(w http.ResponseWriter, resp string) {
	w.WriteHeader(200)
	io.WriteString(w, resp)
}

func sendBadResponse(w http.ResponseWriter, resp string) {
	w.WriteHeader(202)
	io.WriteString(w, resp)
}
