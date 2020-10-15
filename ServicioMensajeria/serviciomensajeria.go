package serviciomensajeria

import (
	"log"

	"golang.org/x/net/context"
)

//Server simple
type Server struct {
	id int
}

//InformaEntrega recibe paquete de Camiones en Logistica
func (s *Server) InformaEntrega(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}
