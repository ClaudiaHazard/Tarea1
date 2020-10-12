package main

import (
	"log"

	"fmt"

	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete/EnviaPaqueteGo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//IP local 10.6.40.161
const (
	ipport = "10.6.40.162/24:50051"
)

func main() {
	fmt.Println("Inicia Camiones")

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	fmt.Println("Crea Conexion")

	c := enviapaquete.NewConexionServiceClient(conn)

	fmt.Println("Envia Mensaje")

	response, err := c.SayHello(context.Background(), &enviapaquete.Message{Body: "Hello From Client!"})

	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)
}
