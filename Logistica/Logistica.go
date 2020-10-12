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
	portCamiones   = "50051"
	ipportCamiones = "10.6.40.162:" + portCamiones
	portCliente    = "50052"
	ipportCliente  = "10.6.40.162:" + portCliente
)

//Para usar en local, cambiar ipportCamiones por ":"+portCamiones y ipportCliente por ":"+portCliente
func main() {
	fmt.Println("Inicia Logistica en espera de mensajes Camiones")
	lis, err := net.Listen("tcp", ipportCamiones)

	if err != nil {
		log.Fatalf("Failed to listen on port "+portCamiones+": %v", err)
	}

	fmt.Println("Crea server para Camiones")
	grpcServerCamion := grpc.NewServer()

	sCamion := enviapaquete.Server{}

	fmt.Println("Crea Conexion de paquetes del Camion")
	enviapaquete.RegisterEnviaPaqueteServiceServer(grpcServerCamion, &sCamion)

	if err := grpcServerCamion.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+portCamiones+": %v", err)
	}

	fmt.Println("Inicia Logistica en espera de mensajes Clientes")
	lis, err2 := net.Listen("tcp", ipportCliente)

	if err2 != nil {
		log.Fatalf("Failed to listen on port "+portCliente+": %v", err2)
	}

	fmt.Println("Crea server para Clientes")
	grpcServerCliente := grpc.NewServer()

	sCliente := enviaorden.Server{}

	fmt.Println("Crea Conexion de ordenes Cliente")
	enviaorden.RegisterEnviaOrdenServiceServer(grpcServerCliente, &sCliente)

	if err2 := grpcServerCliente.Serve(lis); err2 != nil {
		log.Fatalf("Failed to serve gRPC server over "+portCliente+": %v", err2)
	}
}
