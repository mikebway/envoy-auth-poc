package main

import (
	"log"
	"net"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/mikebway/envoy-auth-poc/extauth/server"
	"google.golang.org/grpc"
)

// main is the command line entry point to start the external authorization service.
func main() {

	// List on port 5000
	const address = ":5000"

	// Log that we are starting
	log.Println("Starting server at: " + address)

	// Register this server with Envoy [What does this actually do? The docs are non-existent!]
	s := grpc.NewServer()
	auth.RegisterAuthorizationServer(s, &server.AuthorizationServer{})

	// Set the network listener address
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed setting network address (%s): %v", address, err)
	}

	// Start the service
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed running service: %v", err)
	}
}
