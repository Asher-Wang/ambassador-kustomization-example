package serve

import (
	"context"
	"log"
	"strings"

	"github.com/dizys/ambassador-kustomization-example/auth-service-grpc/config"
	envoy_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	rpc_status "google.golang.org/genproto/googleapis/rpc/status"
)

type AuthServer struct {
}

func (a *AuthServer) Check(context context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	httpReq := req.Attributes.Request.Http
	headers := httpReq.Headers

	authStr, ok := headers["Authorization"]

	if config.Config.GetBool("request_logging") {
		log.Printf("[Request] %s - %s (token: %s): %s\n", httpReq.Method, httpReq.Path+"?"+httpReq.Query, authStr, httpReq.Body)
	}

	if !ok {
		return makeAuthCheckDeniedResponse(int32(rpc.UNAUTHENTICATED), 401, []*envoy_core.HeaderValueOption{}, "Unauthenticated"), nil
	}

	if !strings.HasPrefix(authStr, "Bearer ") {
		return makeAuthCheckDeniedResponse(int32(rpc.PERMISSION_DENIED), 401, []*envoy_core.HeaderValueOption{}, "Invalid access token type"), nil
	}

	unverifiedToken := strings.TrimPrefix(authStr, "Bearer ")
	accessTokens := config.Config.GetStringSlice("access_tokens")

	verified := false

	for _, accessToken := range accessTokens {
		if accessToken == unverifiedToken {
			verified = true
			break
		}
	}

	if !verified {
		return makeAuthCheckDeniedResponse(int32(rpc.PERMISSION_DENIED), 401, []*envoy_core.HeaderValueOption{}, "Unauthorized"), nil
	}

	return makeAuthCheckOkResponse("OK", []*envoy_core.HeaderValueOption{}), nil
}

func makeAuthCheckOkResponse(body string, headers []*envoy_core.HeaderValueOption) *auth.CheckResponse {
	if config.Config.GetBool("request_logging") {
		log.Printf("[Response] %d: %s\n", 200, body)
	}

	return &auth.CheckResponse{
		Status: &rpc_status.Status{
			Code: int32(rpc.OK),
		},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				Headers: headers,
			},
		},
	}
}

func makeAuthCheckDeniedResponse(rpcStatusCode int32, httpStatusCode int32, headers []*envoy_core.HeaderValueOption, body string) *auth.CheckResponse {
	if config.Config.GetBool("request_logging") {
		log.Printf("[Response] %d: %s\n", httpStatusCode, body)
	}

	return &auth.CheckResponse{
		Status: &rpc_status.Status{
			Code: rpcStatusCode,
		},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode(httpStatusCode),
				},
				Headers: headers,
				Body:    body,
			},
		},
	}
}
