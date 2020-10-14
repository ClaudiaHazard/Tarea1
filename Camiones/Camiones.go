package main

import (
	"fmt"
	"log"
	"net"
	"time"

	enviainstrucciones "github.com/ClaudiaHazard/Tarea1/Camiones/EnviaInstrucciones"
	informapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/InformaPaquete"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//IP local 10.6.40.161
const (
	//ipportLogistica = "10.6.40.162:50051"
	ipportLogistica = ":50051"
	//ipportCamiones = "10.6.40.161:50052"
	ipportCamiones = ":50052"
)

//IniciaServidor inicia servidor listen para los servicios correspondientes
func IniciaServidor() {
	lis, err := net.Listen("tcp", ":"+ipportCamiones)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipportCamiones+": %v", err)
	}

	grpcServer := grpc.NewServer()

	sLogistica := enviainstrucciones.Server{}

	fmt.Println("En espera de instrucciones de reparto")
	enviainstrucciones.RegisterEnviaInstruccionesServiceServer(grpcServer, &sLogistica)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipportCamiones+": %v", err)
	}
}

//IniciaCliente inicia conexion cliente
func IniciaCliente() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(ipportLogistica, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return conn
}

//InformaPaqueteLogistica Camion informa estado del paquete a Logistica
func InformaPaqueteLogistica(conn *grpc.ClientConn) string {
	c := informapaquete.NewInformaPaqueteServiceClient(conn)
	response, err := c.InformaPaquete(context.Background(), &informapaquete.Message{Body: "Hola por parte de Camiones!"})

	if err != nil {
		log.Fatalf("Error al llamar InformaPaquete: %s", err)
	}

	log.Printf("Respuesta de Logistica: %s", response.Body)
	return response.Body
}

func main() {
	IniciaServidor()
	var conn = IniciaCliente()

	go InformaPaqueteLogistica(conn)
	go InformaPaqueteLogistica(conn)

	time.Sleep(10 * time.Second)

}
