package handler

import (
	"log"
	"net/http"

	"github.com/dizys/ambassador-kustomization-example/auth-service-http/config"
)

func err(resp http.ResponseWriter, statusCode int, message string) {
	resp.WriteHeader(statusCode)
	resp.Write([]byte(message))

	if config.Config.GetBool("request_logging") {
		log.Printf("[Response] %d: %s\n", statusCode, message)
	}
}
