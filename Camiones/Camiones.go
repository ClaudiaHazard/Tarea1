package main

import (
	"log"
	"time"

	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//IP local 10.6.40.161
const (
	//ipport = "10.6.40.162:50051"
	ipport = ":50051"
)

//EnviaPaquete de Camion a Logistica
func EnviaPaquete(conn *grpc.ClientConn) string {
	c := enviapaquete.NewEnviaPaqueteServiceClient(conn)
	response, err := c.SayHello(context.Background(), &enviapaquete.Message{Body: "Hola por parte de Camiones!"})

	if err != nil {
		log.Fatalf("Error al llamar SayHello: %s", err)
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

	go EnviaPaquete(conn)
	go EnviaPaquete(conn)

	time.Sleep(10 * time.Second)

}
