package servidorcliente

import (
	"fmt"
	"net/http"

	enviaorden "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaOrden"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	portCliente   = "50052"
	ipportCliente = "10.6.40.162:" + portCliente
)

//IniciarServidorCliente Para usar en local, cambiar ipportCamiones por ":"+portCamiones y ipportCliente por ":"+portCliente
func IniciarServidorCliente() {

	fmt.Println("Inicia Logistica en espera de mensajes Clientes")
	//lis, err2 := net.Listen("tcp", ":"+portCliente)

	//if err2 != nil {
	//	log.Fatalf("Failed to listen on port "+portCliente+": %v", err2)
	//}

	fmt.Println("Crea server para Clientes")
	grpcServerCliente := grpc.NewServer()

	sCliente := enviaorden.Server{}

	fmt.Println("Crea Conexion de ordenes Cliente")
	enviaorden.RegisterEnviaOrdenServiceServer(grpcServerCliente, &sCliente)

	http.ListenAndServe(":"+portCliente, grpcServerCliente)

}
