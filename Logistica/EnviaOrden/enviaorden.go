package enviaorden

import (
	"log"

	"golang.org/x/net/context"
)

//Server Test
type Server struct {
}

//EnviaOrden recibe ordenes de los clientes en Logistica
func (s *Server) EnviaOrden(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}
