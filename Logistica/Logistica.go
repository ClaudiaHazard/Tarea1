package main

import (
	"log"
	"net"

	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete/Queestapasando"
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

	grpcServer := grpc.NewServer()

	s := enviapaquete.Server{}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+port+": %v", err)
	}
}
