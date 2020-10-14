package entregaposicion

import (
	"log"

	"golang.org/x/net/context"
)

//Server simple
type Server struct {
}

//EntregaPosicion recibe posicion del camion en Logistica
func (s *Server) EntregaPosicion(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}
