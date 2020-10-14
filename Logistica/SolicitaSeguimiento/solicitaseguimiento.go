package solicitaseguimiento

import (
	"log"

	"golang.org/x/net/context"
)

//Server simple
type Server struct {
}

//SolicitaSeguimiento recibe numero de seguimiento en Logistica
func (s *Server) SolicitaSeguimiento(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}
