package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/dizys/ambassador-kustomization-example/auth-service-http/config"
	"github.com/dizys/ambassador-kustomization-example/auth-service-http/utils"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	authStr := req.Header.Get("Authorization")

	if config.Config.GetBool("request_logging") {
		log.Printf("[Request] %s - %s: %s\n", req.Method, req.RequestURI, req.PostForm.Encode())
	}

	if authStr == "" {
		utils.Err(resp, 401, "Unauthenticated")
	}

	if !strings.HasPrefix(authStr, "Bearer ") {
		utils.Err(resp, 401, "Invalid access token type")
		return
	}

	unverifiedToken := strings.TrimPrefix(authStr, "Bearer ")
	accessTokens := config.Config.GetStringSlice("accessTokens")

	verified := false

	for _, accessToken := range accessTokens {
		if accessToken == unverifiedToken {
			verified = true
			break
		}
	}

	if !verified {
		utils.Err(resp, 401, "Unauthorized")
		return
	}

	resp.Write([]byte("OK"))

	if config.Config.GetBool("request_logging") {
		log.Printf("[Response] 200: OK\n")
	}
}
