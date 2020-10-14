package muestraseguimiento

import (
	"log"

	"golang.org/x/net/context"
)

//Server simple
type Server struct {
}

//MuestraSeguimiento recibe estado del paquete en Cliente
func (s *Server) MuestraSeguimiento(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}
