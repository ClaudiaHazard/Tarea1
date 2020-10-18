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

var clienteid string
var camion int
var arrRetail []*sm.Paquete
var arrPrioritario []*sm.Paquete
var arrNormal []*sm.Paquete

//CodSeg Codigo de seguimiento que se incrementa en uno cada vez que se genera un nuevo codigo
var CodSeg int32

//Server datos
type Server struct {
	clienteid      string
	arrRetail      []*sm.Paquete
	arrPrioritario []*sm.Paquete
	arrNormal      []*sm.Paquete
}

//AgregaACola agrega paquete a cola correspondiente
func AgregaACola(p *sm.Paquete, s *Server) {

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
func BorrarElemento(arr []*sm.Paquete, pos int) []*sm.Paquete {
	copy(arr[pos:], arr[pos+1:])    // Shift a[i+1:] left one index.
	arr[len(arr)-1] = &sm.Paquete{} // Erase last element (write zero value).
	arr = arr[:len(arr)-1]
	return arr
}

//AsignaPaquete asigna paquete al tipo de camion correspondiente.
func AsignaPaquete(s *Server, tipoCam string, entrPrevRetail bool, paqCargRetail bool) *sm.Paquete {
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
			return &sm.Paquete{}
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
			return &sm.Paquete{}
		}
	}
	return &sm.Paquete{}
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
	return &sm.Message{Body: "Hola desde Logistica! camion numero " + s.clienteid}, nil
}

//RecibeInstrucciones Camion avisa que esta disponible y se le envia paquete
func (s *Server) RecibeInstrucciones(ctx context.Context, in *sm.DisponibleCamion) (*sm.Paquete, error) {
	log.Printf("El Camion %d se encuentra disponible.", in.Id)
	paq := AsignaPaquete(s, in.Tipo, in.EntrPrevRetail, in.PaqCargRetail)
	return paq, nil
}

//RealizaOrden cliente envia orden, logistica retorna Codigo de seguimiento
func (s *Server) RealizaOrden(ctx context.Context, in *sm.Orden) (*sm.CodSeguimiento, error) {
	log.Printf("Se recibio paquete %s con Id: %s", in.Nombre, in.Id)

	paq := CreaPaquete(in)
	AgregaACola(paq, s)

	return &sm.CodSeguimiento{CodigoSeguimiento: paq.CodigoSeguimiento}, nil

}

//SolicitaSeguimiento solicita estado de su orden
func (s *Server) SolicitaSeguimiento(ctx context.Context, in *sm.CodSeguimiento) (*sm.Estado, error) {
	log.Printf("Receive message body from client: %d y %s", in.CodigoSeguimiento, s.clienteid)
	return &sm.Estado{Estado: "Bonito"}, nil
}

//CreaPaquete genera paquete de la orden que entrego el Cliente
func CreaPaquete(o *sm.Orden) *sm.Paquete {
	if o.Tipo == "Normal" || o.Tipo == "Prioritario" {
		CodSeg = CodSeg + 1
		return &sm.Paquete{Id: o.Id, CodigoSeguimiento: CodSeg, Tipo: o.Tipo, Valor: o.Valor, Intentos: 0, Estado: "En bodega", Origen: o.Origen, Destino: o.Destino, Nombre: o.Nombre}
	}
	return &sm.Paquete{Id: o.Id, CodigoSeguimiento: 0, Tipo: o.Tipo, Valor: o.Valor, Intentos: 0, Estado: "En bodega", Origen: o.Origen, Destino: o.Destino, Nombre: o.Nombre}
}

//Para usar en local, cambiar ipport por ":"+port
func main() {
	lis, err := net.Listen("tcp", ipport)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipport+": %v", err)
	}

	s := Server{"1", []*sm.Paquete{}, []*sm.Paquete{}, []*sm.Paquete{}}

	CodSeg = 0

	paq := &sm.Paquete{Id: "1", CodigoSeguimiento: 1, Tipo: "Retail", Valor: 10, Intentos: 0, Estado: "En bodega", Origen: "Origen A", Destino: "Destino A", Nombre: "Bicicleta"}

	s.arrRetail = append(s.arrRetail, paq)

	grpcServer := grpc.NewServer()

	fmt.Println("En espera de Informacion paquetes para servidor")

	sm.RegisterMensajeriaServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipport+": %v", err)
	}

}
