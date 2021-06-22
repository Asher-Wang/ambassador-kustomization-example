package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dizys/ambassador-kustomization-example/auth-service-grpc/config"
	"github.com/dizys/ambassador-kustomization-example/auth-service-grpc/serve"
)

func main() {
	config.SetupConfig()

	port := config.Config.GetInt("port")

	if port <= 0 || port >= 65_535 {
		log.Fatalf("Invalid port number: %d\n", port)
		os.Exit(1)
	}

	s, err := serve.CreateServer()

	if err != nil {
		log.Fatalf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("Error trying to listen on tcp port %d: %v\n", port, err)
		os.Exit(1)
	}

	log.Printf("Auth service (gRPC) running on %d...\n", port)
	s.Serve(listener)
}
