package main

import (
	"fmt"
	"log"
	"net"

	serviciomensajeria "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	//ipportLogistica = "10.6.40.162:50051"
	ipport = ":50051"
)

//IniciaServidor inicia servidor listen para los servicios correspondientes
func IniciaServidor() {
	lis, err := net.Listen("tcp", ipport)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipport+": %v", err)
	}

	grpcServer := grpc.NewServer()

	s := serviciomensajeria.Server{id: 2}

	fmt.Println("En espera de Informacion paquetes")

	serviciomensajeria.RegisterMensajeriaServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipport+": %v", err)
	}
}

//Para usar en local, cambiar ipport por ":"+port
func main() {
	IniciaServidor()

}
