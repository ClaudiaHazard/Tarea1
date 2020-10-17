package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	sm "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	//ipportLogistica = "10.6.40.162:50051"
	ipport = ":50051"
)

//Server datos locales de Logistica
type Server struct {
	camion         int
	arrRetail      []sm.Paquete
	arrPrioritario []sm.Paquete
	arrNormal      []sm.Paquete
}

//AgregaACola agrega paquete a cola correspondiente
func AgregaACola(p sm.Paquete, s Server) {
	if p.Tipo == "Retail" {
		s.arrRetail = append(s.arrRetail, p)
	}
	if p.Tipo == "Prioritario" {
		s.arrPrioritario = append(s.arrPrioritario, p)
	}
	if p.Tipo == "Retail" {
		s.arrNormal = append(s.arrNormal, p)
	}
}

//BorrarElemento borra el elemento en la posicion pos.
func BorrarElemento(arr []sm.Paquete, pos int) []sm.Paquete {
	copy(arr[pos:], arr[pos+1:])   // Shift a[i+1:] left one index.
	arr[len(arr)-1] = sm.Paquete{} // Erase last element (write zero value).
	arr = arr[:len(arr)-1]
	return arr
}

//AsignaPaquete asigna paquete al tipo de camion correspondiente.
func AsignaPaquete(s *Server, tipoCam string, entrPrevRetail bool, paqCargRetail bool) sm.Paquete {
	if tipoCam == "Normal" {
		if len(s.arrPrioritario) != 0 {
			p := s.arrPrioritario[0]
			s.arrPrioritario = BorrarElemento(s.arrPrioritario, 0)
			return p
		} else if len(s.arrNormal) != 0 {
			p := s.arrNormal[0]
			s.arrNormal = BorrarElemento(s.arrNormal, 0)
			return p
		} else {
			return sm.Paquete{}
		}
	}
	if tipoCam == "Retail" {
		if len(s.arrRetail) != 0 {
			p := s.arrRetail[0]
			s.arrRetail = BorrarElemento(s.arrRetail, 0)
			return p
		} else if len(s.arrPrioritario) != 0 && entrPrevRetail && paqCargRetail {
			p := s.arrPrioritario[0]
			s.arrPrioritario = BorrarElemento(s.arrPrioritario, 0)
			return p
		} else {
			return sm.Paquete{}
		}
	}
	return sm.Paquete{}
}

//EntregaPosicion recibe paquete de Camiones en Logistica
func (s *Server) EntregaPosicion(ctx context.Context, in *sm.InformacionPaquete) (*sm.Message, error) {
	log.Printf("Receive message body from client: %d", in.CodigoSeguimiento)
	//s = ctx.Value
	return &sm.Message{Body: "Hola desde Logistica!"}, nil
}

//InformaEntrega recibe paquete de Camiones en Logistica
func (s *Server) InformaEntrega(ctx context.Context, in *sm.Message) (*sm.Message, error) {
	log.Printf("Receive message body from client: %s yep %d", in.Body, s.camion)
	return &sm.Message{Body: "Hola desde Logistica! camion numero " + strconv.Itoa(s.camion)}, nil
}

//RecibeInstrucciones recibe paquete de Camiones en Logistica
func (s *Server) RecibeInstrucciones(ctx context.Context, in *sm.Message) (*sm.Paquete, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &sm.Paquete{}, nil
}

//RealizaOrden recibe paquete de Camiones en Logistica
func (s *Server) RealizaOrden(ctx context.Context, in *sm.Orden) (*sm.CodSeguimiento, error) {
	log.Printf("Receive message body from client: %s", in.Nombre)
	return &sm.CodSeguimiento{}, nil
}

//SolicitaSeguimiento recibe paquete de Camiones en Logistica
func (s *Server) SolicitaSeguimiento(ctx context.Context, in *sm.CodSeguimiento) (*sm.Estado, error) {
	log.Printf("Receive message body from client: %d y %d", in.CodigoSeguimiento, s.camion)
	return &sm.Estado{Estado: "Bonito"}, nil
}

//Para usar en local, cambiar ipport por ":"+port
func main() {
	lis, err := net.Listen("tcp", ipport)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipport+": %v", err)
	}

	s := Server{}

	grpcServer := grpc.NewServer()

	fmt.Println("En espera de Informacion paquetes para servidor")

	sm.RegisterMensajeriaServiceServer(grpcServer, &s)

	s = Server{}

	if s.camion == 3 {
		lis.Close()
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipport+": %v", err)
	}

}
