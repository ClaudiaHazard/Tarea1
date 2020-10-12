package enviapaquete

import (
	"log"

	"golang.org/x/net/context"
)

//Server Test
type Server struct {
}

//SayHello Test
func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}
