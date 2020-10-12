package servidorcliente

import (
	"fmt"
	"log"
	"net"

	enviaorden "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaOrden"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	portCliente   = "50052"
	ipportCliente = "10.6.40.162:" + portCliente
)

//Para usar en local, cambiar ipportCamiones por ":"+portCamiones y ipportCliente por ":"+portCliente
func servidorcliente() {

	fmt.Println("Inicia Logistica en espera de mensajes Clientes")
	lis, err2 := net.Listen("tcp", ":"+portCliente)

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
