package main

import (
	"fmt"
	"log"
	"net"

	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete/EnviaPaqueteGo"
	"google.golang.org/grpc"
)

const (
	port = "50051"
	//ipport = "10.6.40.161/24:" + port
	ipport = "10.6.40.162/24:" + port
)

func main() {
	fmt.Println("Inicia Logistica")
	lis, err := net.Listen("tcp", ipport)

	if err != nil {
		log.Fatalf("Failed to listen on port "+port+": %v", err)
	}

	fmt.Println("Crea server")
	grpcServer := grpc.NewServer()

	s := enviapaquete.Server{}

	fmt.Println("Crea Conexion")
	enviapaquete.RegisterConexionServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+port+": %v", err)
	}
}
