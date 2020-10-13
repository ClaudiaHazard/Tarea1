package main

import (
	"fmt"
	"log"
	"net"

	enviaorden "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaOrden"
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

	//lis, err := net.Listen("tcp", ":"+portCamiones)
	lis2, err2 := net.Listen("tcp", ":"+portCliente)

	//if err != nil {
	//	log.Fatalf("Failed to listen on port "+portCamiones+": %v", err)
	//}

	if err2 != nil {
		log.Fatalf("Failed to listen on port "+portCliente+": %v", err2)
	}

	//grpcServerCamion := grpc.NewServer()
	grpcServerCliente := grpc.NewServer()

	//sCamion := enviapaquete.Server{}
	sCliente := enviaorden.Server{}

	//fmt.Println("Envia paquetes")
	//enviapaquete.RegisterEnviaPaqueteServiceServer(grpcServerCamion, &sCamion)
	fmt.Println("Envia ordenes")
	enviaorden.RegisterEnviaOrdenServiceServer(grpcServerCliente, &sCliente)

	//if err := grpcServerCamion.Serve(lis); err != nil {
	//	log.Fatalf("Failed to serve gRPC server over "+portCamiones+": %v", err)
	//}

	if err2 := grpcServerCliente.Serve(lis2); err2 != nil {
		log.Fatalf("Failed to serve gRPC server over "+portCliente+": %v", err2)
	}
}
