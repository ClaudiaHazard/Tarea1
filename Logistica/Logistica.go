package main

import (
	"fmt"
	"log"
	"net"

	enviaorden "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaOrden"
	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	port   = "50051"
	ipport = "10.6.40.162:" + port
)

//Para usar en local, cambiar ipport por ":"+port
func main() {

	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Failed to listen on port "+port+": %v", err)
	}

	grpcServer := grpc.NewServer()

	sCamion := enviapaquete.Server{}
	sCliente := enviaorden.Server{}

	fmt.Println("Envia paquetes")
	enviapaquete.RegisterEnviaPaqueteServiceServer(grpcServer, &sCamion)
	fmt.Println("Envia ordenes")
	enviaorden.RegisterEnviaOrdenServiceServer(grpcServer, &sCliente)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+port+": %v", err)
	}
}
