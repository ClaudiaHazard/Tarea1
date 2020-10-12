package servidorcamiones

import (
	"fmt"
	"log"
	"net"

	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	portCamiones   = "50051"
	ipportCamiones = "10.6.40.162:" + portCamiones
)

//IniciarServidorCamiones Para usar en local, cambiar ipportCamiones por ":"+portCamiones y ipportCliente por ":"+portCliente
func IniciarServidorCamiones() {
	fmt.Println("Inicia Logistica en espera de mensajes Camiones")
	lis, err := net.Listen("tcp", ":"+portCamiones)

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

}
