package serve

import (
	"fmt"

	"github.com/dizys/ambassador-kustomization-example/auth-service-grpc/config"
	"google.golang.org/grpc"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

func CreateServer() (*grpc.Server, error) {
	maxConcurrentStreams := config.Config.GetUint32("max_concurrent_streams")

	if maxConcurrentStreams <= 0 {
		return nil, fmt.Errorf("invalid max concurrent stream number: %d", maxConcurrentStreams)
	}

	opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(maxConcurrentStreams)}

	server := grpc.NewServer(opts...)

	auth.RegisterAuthorizationServer(server, &AuthServer{})

	return server, nil
}
