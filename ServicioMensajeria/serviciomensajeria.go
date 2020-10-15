package serviciomensajeria

import (
	"log"

	"golang.org/x/net/context"
)

//Server simple
type Server struct {
	//id int
}

//EntregaPosicion recibe paquete de Camiones en Logistica
func (s *Server) EntregaPosicion(ctx context.Context, in *InformacionPaquete) (*Message, error) {
	log.Printf("Receive message body from client: %d", in.CodigoSeguimiento)
	return &Message{Body: "Hola desde Logistica!"}, nil
}

//InformaEntrega recibe paquete de Camiones en Logistica
func (s *Server) InformaEntrega(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hola desde Logistica!"}, nil
}

//RecibeInstrucciones recibe paquete de Camiones en Logistica
func (s *Server) RecibeInstrucciones(ctx context.Context, in *Message) (*Paquete, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Paquete{}, nil
}

//RealizaOrden recibe paquete de Camiones en Logistica
func (s *Server) RealizaOrden(ctx context.Context, in *Orden) (*CodSeguimiento, error) {
	log.Printf("Receive message body from client: %s", in.Nombre)
	return &CodSeguimiento{}, nil
}

//SolicitaSeguimiento recibe paquete de Camiones en Logistica
func (s *Server) SolicitaSeguimiento(ctx context.Context, in *CodSeguimiento) (*Estado, error) {
	log.Printf("Receive message body from client: %d", in.CodigoSeguimiento)
	return &Estado{Estado: "Bonito"}, nil
}
