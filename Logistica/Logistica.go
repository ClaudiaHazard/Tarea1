package main

import (
	"fmt"
	"log"
	"net"

	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete/EnviaPaqueteGo"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	port   = "50051"
	ipport = "10.6.40.162:" + port
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

	fmt.Println("Funciono lo creado")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+port+": %v", err)
	}
}
