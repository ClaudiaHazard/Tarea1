package main

import (
	"context"
	"fmt"
	"log"
	"net"

	enviainstrucciones "github.com/ClaudiaHazard/Tarea1/Camiones/EnviaInstrucciones"
	enviaorden "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaOrden"
	informapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/InformaPaquete"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	//ipportLogistica = "10.6.40.162:50051"
	ipportLogistica = ":50051"
	//ipportCamiones = "10.6.40.161:50051"
	ipportCamiones = ":50051"
)

//IniciaServidor inicia servidor listen para los servicios correspondientes
func IniciaServidor() {
	lis, err := net.Listen("tcp", ipportLogistica)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipportLogistica+": %v", err)
	}

	grpcServer := grpc.NewServer()

	sCamion := informapaquete.Server{1}
	sCliente := enviaorden.Server{}

	fmt.Println("En espera de Informacion paquetes")
	informapaquete.RegisterInformaPaqueteServiceServer(grpcServer, &sCamion)
	fmt.Println("En espera de nuevas ordenes de cliente")
	enviaorden.RegisterEnviaOrdenServiceServer(grpcServer, &sCliente)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipportLogistica+": %v", err)
	}
}

//IniciaCliente inicia conexion cliente
func IniciaCliente() *grpc.ClientConn {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(ipportCamiones, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return conn
}

//EnviaInstrucciones de Camion a Logistica
func EnviaInstrucciones(conn *grpc.ClientConn) string {
	c := enviainstrucciones.NewEnviaInstruccionesServiceClient(conn)
	response, err := c.EnviaInstrucciones(context.Background(), &enviainstrucciones.Message{Body: "Hola por parte de Logistica!"})

	if err != nil {
		log.Fatalf("Error al llamar InformaPaquete: %s", err)
	}

	log.Printf("Respuesta de Logistica: %s", response.Body)
	return response.Body
}

//Para usar en local, cambiar ipport por ":"+port
func main() {
	go IniciaServidor()

}
