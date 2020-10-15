package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	//ipportLogistica = "10.6.40.162:50051"
	ipport = ":50051"
)

//IniciaServidor inicia servidor listen para los servicios correspondientes
func IniciaServidor() {
	lis, err := net.Listen("tcp", ipport)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipport+": %v", err)
	}

	grpcServer := grpc.NewServer()

	s := Server{}

	fmt.Println("En espera de Informacion paquetes")

	sm.RegisterMensajeriaServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipport+": %v", err)
	}
}

//Server simple
type Server struct {
	//id int
}

//EntregaPosicion recibe paquete de Camiones en Logistica
func (s *Server) EntregaPosicion(ctx context.Context, in *InformacionPaquete) (*logistica.Message, error) {
	log.Printf("Receive message body from client: %d", in.CodigoSeguimiento)
	return &Message{Body: "Hola desde Logistica!"}, nil
}

//InformaEntrega recibe paquete de Camiones en Logistica
func (s *Server) InformaEntrega(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}

//RecibeInstrucciones recibe paquete de Camiones en Logistica
func (s *Server) RecibeInstrucciones(ctx context.Context, in *Message) (*Paquete, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Paquete{}, nil
}

//RealizaOrden recibe paquete de Camiones en Logistica
func (s *Server) RealizaOrden(ctx context.Context, in *Orden) (*CodSeguimiento, error) {
	log.Printf("Receive message body from client: %s", in.Nombre)
	return &CodSeguimiento{}, nil
}

//SolicitaSeguimiento recibe paquete de Camiones en Logistica
func (s *Server) SolicitaSeguimiento(ctx context.Context, in *CodSeguimiento) (*Estado, error) {
	log.Printf("Receive message body from client: %d", in.CodigoSeguimiento)
	return &Estado{Estado: "Bonito"}, nil
}

//Para usar en local, cambiar ipport por ":"+port
func main() {
	lis, err := net.Listen("tcp", ipport)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipport+": %v", err)
	}

	grpcServer := grpc.NewServer()

	s := Server{}

	fmt.Println("En espera de Informacion paquetes")

	RegisterMensajeriaServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipport+": %v", err)
	}

}
