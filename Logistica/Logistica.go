package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = "50051"
)

func main() {
	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Failed to listen on port "+port+": %v", err)
	}

	//s := enviapaquete.Server{}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+port+": %v", err)
	}
}
