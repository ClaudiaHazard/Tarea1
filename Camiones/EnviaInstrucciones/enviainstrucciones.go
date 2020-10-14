package enviainstrucciones

import (
	"log"

	"golang.org/x/net/context"
)

//Server simple
type Server struct {
}

//EnviaInstrucciones recibe instrucciones del paquete en Camiones
func (s *Server) EnviaInstrucciones(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}
