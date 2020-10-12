package main

import (
	"log"

	"fmt"

	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete/EnviaPaqueteGo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = "50051"
)

func main() {
	fmt.Println("Helloww world")

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := enviapaquete.NewConexionServiceClient(conn)

	response, err := c.SayHello(context.Background(), &enviapaquete.Message{Body: "Hello From Client!"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Body)
}
