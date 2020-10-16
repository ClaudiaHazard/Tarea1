package main

import (
	"log"
	"time"

	serviciomensajeria "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//IP local 10.6.40.161
const (
	//ipport = "10.6.40.162:50051"
	ipport = ":50051"
)

//Camiones tipo camion 1,2,3.
type Camiones struct {
	m map[string]string
}

//Get Obtiene el camion segun numero camion
func (cam Camiones) Get(numeroCamion string) string {
	return cam.m[numeroCamion]
}

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
	cam := Camiones{map[string]string{
		"1": "one",
		"2": "two",
	}}
	ctx := context.Background()
	response, err := c.InformaEntrega(context.WithValue(ctx, "camiones", cam), &serviciomensajeria.Message{Body: "Hola por parte de Camiones!"})

	if err != nil {
		log.Fatalf("Error al llamar InformaPaquete: %s", err)
	}

	log.Printf("Respuesta de Logistica: %s", response.Body)
	return response.Body
}

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	go InformaPaqueteLogistica(conn)
	go InformaPaqueteLogistica(conn)

	time.Sleep(10 * time.Second)

}
