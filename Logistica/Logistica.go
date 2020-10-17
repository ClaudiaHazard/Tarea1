package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	sm "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	//ipportLogistica = "10.6.40.162:50051"
	ipport = ":50051"
)

//Server simple
type Server struct {
	id int
}

//CamionKey key del camion
type CamionKey struct {
	s string
}

//EntregaPosicion recibe paquete de Camiones en Logistica
func (s *Server) EntregaPosicion(ctx context.Context, in *sm.InformacionPaquete) (*sm.Message, error) {
	log.Printf("Receive message body from client: %d", in.CodigoSeguimiento)
	s = ctx.Value
	return &sm.Message{Body: "Hola desde Logistica!"}, nil
}

//InformaEntrega recibe paquete de Camiones en Logistica
func (s *Server) InformaEntrega(ctx context.Context, in *sm.Message) (*sm.Message, error) {
	log.Printf("Receive message body from client: %s yep %d", in.Body, s.id)
	return &sm.Message{Body: "Hola desde Logistica! camion numero " + strconv.Itoa(s.id)}, nil
}

//RecibeInstrucciones recibe paquete de Camiones en Logistica
func (s *Server) RecibeInstrucciones(ctx context.Context, in *sm.Message) (*sm.Paquete, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &sm.Paquete{}, nil
}

//RealizaOrden recibe paquete de Camiones en Logistica
func (s *Server) RealizaOrden(ctx context.Context, in *sm.Orden) (*sm.CodSeguimiento, error) {
	log.Printf("Receive message body from client: %s", in.Nombre)
	return &sm.CodSeguimiento{}, nil
}

//SolicitaSeguimiento recibe paquete de Camiones en Logistica
func (s *Server) SolicitaSeguimiento(ctx context.Context, in *sm.CodSeguimiento) (*sm.Estado, error) {
	log.Printf("Receive message body from client: %d y %d", in.CodigoSeguimiento, s.id)
	return &sm.Estado{Estado: "Bonito"}, nil
}

//Para usar en local, cambiar ipport por ":"+port
func main() {
	lis, err := net.Listen("tcp", ipport)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipport+": %v", err)
	}

	s := Server{1}

	grpcServer := grpc.NewServer()

	fmt.Println("En espera de Informacion paquetes para servidor")

	sm.RegisterMensajeriaServiceServer(grpcServer, &s)

	s = Server{2}

	if s.id == 3 {
		lis.Close()
	}

	s = Server{3}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipport+": %v", err)
	}

}
