package main

import (
	"log"
	"time"

	serviciocamiones "github.com/ClaudiaHazard/Tarea1/ServiciosServidor/ServicioCamiones"
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
	c := serviciocamiones.NewInformaPaqueteServiceClient(conn)
	response, err := c.InformaPaquete(context.Background(), &informapaquete.Message{Body: "Hola por parte de Camiones!"})

	if err != nil {
		log.Fatalf("Error al llamar InformaPaquete: %s", err)
	}

	log.Printf("Respuesta de Logistica: %s", response.Body)
	return response.Body
}

func main() {

	var conn = IniciaCliente()

	go InformaPaqueteLogistica(conn)
	go InformaPaqueteLogistica(conn)
	time.Sleep(10 * time.Second)

}
