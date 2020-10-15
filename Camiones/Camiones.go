package main

import (
	"log"

	serviciomensajeria "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//IP local 10.6.40.161
const (
	//ipport = "10.6.40.162:50051"
	ipport = ":50051"
)

//IniciaCliente inicia conexion cliente
func IniciaCliente() *grpc.ClientConn {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return conn
}

//InformaPaqueteLogistica Camion informa estado del paquete a Logistica
func InformaPaqueteLogistica(conn *grpc.ClientConn) string {
	c := serviciomensajeria.NewMensajeriaServiceClient(conn)
	response, err := c.InformaEntrega(context.Background(), &serviciomensajeria.Message{Body: "Hola por parte de Camiones!"})

	if err != nil {
		log.Fatalf("Error al llamar InformaPaquete: %s", err)
	}

	log.Printf("Respuesta de Logistica: %s", response.Body)
	return response.Body
}

func main() {
	log.Printf("Respuesta de Logisticadsaasdasd")
	var conn = IniciaCliente()
	log.Printf("Respuesta de Logistsdada")
	InformaPaqueteLogistica(conn)

}
