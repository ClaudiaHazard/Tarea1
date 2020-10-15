package main

import (
	"fmt"
	"log"
	"net"

	serviciosservidor "github.com/ClaudiaHazard/Tarea1/ServiciosServidor"
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

	sCamion := serviciosservidor.RegisterCamionesServiceServer.Server{}
	//sCliente := enviaorden.Server{}

	fmt.Println("En espera de Informacion paquetes")
	serviciosservidor.RegisterCamionesServiceServer(grpcServer, &sCamion)
	//fmt.Println("En espera de nuevas ordenes de cliente")
	//enviaorden.RegisterEnviaOrdenServiceServer(grpcServer, &sCliente)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipport+": %v", err)
	}
}

//Para usar en local, cambiar ipport por ":"+port
func main() {
	IniciaServidor()

}
