package conexion

import(
	"log"
	"golang.org/x/net/context"
)

type Server struct {

}

func (s * Server) SayHello(ctx contex.Context, message * Message) (* Message, error ) {
	log.Printf("Received message body from client: %s", message.Body)
	return &Message{Body: "Hello from the Server!"}, nil
}