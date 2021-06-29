package main

import (
	"fmt"
	"log"
	"net/http"
)

var port = 3000

type Handler struct {
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	passport := req.Header.Get("x-passport")

	if passport == "" {
		rErr(resp, 403, "Empty passport")
		return
	}

	resp.Write([]byte(passport))
	resp.Header().Add("content-type", "application/json")
}

func main() {
	log.Printf("Auth service (HTTP) running on %d...\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &Handler{})
	if err != nil {
		panic(err)
	}
}

func rErr(resp http.ResponseWriter, statusCode int, message string) {
	resp.WriteHeader(statusCode)
	resp.Write([]byte(message))

	log.Printf("[Response] %d: %s\n", statusCode, message)
}
