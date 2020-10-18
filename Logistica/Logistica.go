package main

import (
	"context"
	"fmt"
	"log"
	"net"

	sm "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	//ipportLogistica = "10.6.40.162:50051"
	ipport = ":50051"
)

//CamionResp para reconocer llamadas de camion en Logistica.
type CamionResp struct {
	id   int32
	tipo string
}

//Logistica datos locales de Logistica
type Logistica struct {
	camion         int
	arrRetail      []sm.Paquete
	arrPrioritario []sm.Paquete
	arrNormal      []sm.Paquete
}

//Server datos
type Server struct {
	clienteid string
}

//AgregaACola agrega paquete a cola correspondiente
func AgregaACola(p sm.Paquete, s Logistica) {
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
func AsignaPaquete(s *Logistica, tipoCam string, entrPrevRetail bool, paqCargRetail bool) sm.Paquete {
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

//EntregaPosicion Entrega actualizacion de paquete
func (s *Server) EntregaPosicion(ctx context.Context, in *sm.InformacionPaquete) (*sm.Message, error) {
	log.Printf("Receive message body from client: %d", in.Idpaquete)
	//s = ctx.Value
	return &sm.Message{Body: "Hola desde Logistica!"}, nil
}

//InformaEntrega Informaque camion termino orden
func (s *Server) InformaEntrega(ctx context.Context, in *sm.InformePaquetes) (*sm.Message, error) {
	log.Printf("Receive message body from client: yep")
	tipoCam := ctx.Value("tipo")
	log.Printf(tipoCam.(string))
	return &sm.Message{Body: "Hola desde Logistica! camion numero " + s.clienteid}, nil
}

//RecibeInstrucciones Camion avisa que esta disponible y se le envia paquete
func (s *Server) RecibeInstrucciones(ctx context.Context, in *sm.DisponibleCamion) (*sm.Paquete, error) {
	log.Printf("Receive message body from client: %d", in.Id)
	return &sm.Paquete{}, nil
}

//RealizaOrden cliente envia orden
func (s *Server) RealizaOrden(ctx context.Context, in *sm.Orden) (*sm.CodSeguimiento, error) {
	log.Printf("Receive message body from client: %s", in.Nombre)
	return &sm.CodSeguimiento{}, nil
}

//SolicitaSeguimiento solicita estado de su orden
func (s *Server) SolicitaSeguimiento(ctx context.Context, in *sm.CodSeguimiento) (*sm.Estado, error) {
	log.Printf("Receive message body from client: %d y %s", in.CodigoSeguimiento, s.clienteid)
	return &sm.Estado{Estado: "Bonito"}, nil
}

//Para usar en local, cambiar ipport por ":"+port
func main() {
	lis, err := net.Listen("tcp", ipport)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipport+": %v", err)
	}

	s := Server{"Cliente1"}

	grpcServer := grpc.NewServer()

	fmt.Println("En espera de Informacion paquetes para servidor")

	sm.RegisterMensajeriaServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipport+": %v", err)
	}

}
