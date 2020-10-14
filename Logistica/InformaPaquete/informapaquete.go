package informapaquete

import (
	"log"

	"golang.org/x/net/context"
)

//Server simple
type Server struct {
	id int
}

//InformaPaquete recibe paquete de Camiones en Logistica
func (s *Server) InformaPaquete(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}
