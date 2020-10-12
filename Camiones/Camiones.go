package main

import (
	"log"

	"fmt"

	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete/EnviaPaqueteGo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

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

	c := enviapaquete.NewConexionServiceClient(conn)

	fmt.Println("Crea Conexion")

	response, err := c.SayHello(context.Background(), &enviapaquete.Message{Body: "Hello From Client!"})

	fmt.Println("Envia Mensaje")

	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)
}
