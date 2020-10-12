package enviapaquete

import (
	"log"

	enviapaquete "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaPaquete"
	"golang.org/x/net/context"
)

//Server Test
type Server struct {
	enviapaquete.UnimplementedConexionServiceServer
}

//SayHello Test
func (s *Server) SayHello(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received message body from client: %s", message.Body)
	return &Message{Body: "Hello from the Server"}, nil
}
