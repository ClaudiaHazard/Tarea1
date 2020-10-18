package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"

	sm "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"google.golang.org/grpc"
)

//IP local 10.6.40.162
const (
	ipportgrpc = "10.6.40.162:50051"
	//ipportgrpc = ":50051"
	ipportrabbitmq = "amqp://test:test@10.6.40.1:5672/"
	//ipportrabbitmq = "amqp://guest:guest@localhost:5672/"
)

//error handling
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//CamionResp para reconocer llamadas de camion en Logistica.
type CamionResp struct {
	id   int32
	tipo string
}

//CodSeg Codigo de seguimiento que se incrementa en uno cada vez que se genera un nuevo codigo
var CodSeg int32
var conn *amqp.Connection
var err error
var csvFile *os.File

//Server datos
type Server struct {
	clienteid      string
	arrRetail      []*sm.Paquete
	arrPrioritario []*sm.Paquete
	arrNormal      []*sm.Paquete
	Seguimiento    map[int32]string
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
	log.Printf("Recibido estado con codigo de seguimiento: %d", in.CodigoSeguimiento)
	s.Seguimiento[in.CodigoSeguimiento] = in.Estado
	return &sm.Message{Body: "Ok"}, nil
}

//ReporteFinanzas envía a finanzas datos de paquetes completados
func ReporteFinanzas(pa *sm.Paquete, pa2 *sm.Paquete, conn *amqp.Connection) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")

	var entre string
	if pa.Estado == "Recibido" {
		entre = "true"
	} else if pa.Estado == "No Recibido" {
		entre = "false"
	}

	var entre2 string
	if pa2.Estado == "Recibido" {
		entre2 = "true"
	} else if pa2.Estado == "No Recibido" {
		entre2 = "false"
	}

	body := `{"ID":` + `"` + pa.Id + `"` + `, "intentos" :` + strconv.Itoa(int(pa.Intentos)) + `, "entregado":` + entre + `, "valor" :  ` + strconv.Itoa(int(pa.Valor)) + `, "tipo": ` + `"` + pa.Tipo + `"` + `}`
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	//log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")

	body2 := `{"ID":` + `"` + pa2.Id + `"` + `, "intentos" :` + strconv.Itoa(int(pa2.Intentos)) + `, "entregado":` + entre2 + `, "valor" :  ` + strconv.Itoa(int(pa2.Valor)) + `, "tipo": ` + `"` + pa2.Tipo + `"` + `}`
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body2),
		})
	//log.Printf(" [x] Sent %s", body2)
	failOnError(err, "Failed to publish a message")

}

//InformaEntrega Informaque camion termino entrega
func (s *Server) InformaEntrega(ctx context.Context, in *sm.InformePaquetes) (*sm.Message, error) {

	log.Printf("Entrega completada.")
	pa := in.Paquetes[0]

	pa2 := in.Paquetes[1]

	//Aqui se debe enviar mensaje a Finanzas con los 2 paquetes para que calcule lo que deba calcular.

	ReporteFinanzas(pa, pa2, conn)

	return &sm.Message{Body: "Ok"}, nil
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

	//Agrega datos de la orden al registro
	EditaResigtro(csvFile, in, paq.CodigoSeguimiento)

	return &sm.CodSeguimiento{CodigoSeguimiento: paq.CodigoSeguimiento}, nil

}

//SolicitaSeguimiento solicita estado de su orden
func (s *Server) SolicitaSeguimiento(ctx context.Context, in *sm.CodSeguimiento) (*sm.Estado, error) {
	log.Printf("Se envia el estado del paquete con codigode seguimiento: %d", in.CodigoSeguimiento)
	return &sm.Estado{Estado: s.Seguimiento[in.CodigoSeguimiento]}, nil
}

//CreaPaquete genera paquete de la orden que entrego el Cliente
func CreaPaquete(o *sm.Orden) *sm.Paquete {
	if o.Tipo == "Normal" || o.Tipo == "Prioritario" {
		CodSeg = CodSeg + 1
		return &sm.Paquete{Id: o.Id, CodigoSeguimiento: CodSeg, Tipo: o.Tipo, Valor: o.Valor, Intentos: 0, Estado: "En bodega", Origen: o.Origen, Destino: o.Destino, Nombre: o.Nombre}
	}
	return &sm.Paquete{Id: o.Id, CodigoSeguimiento: 0, Tipo: o.Tipo, Valor: o.Valor, Intentos: 0, Estado: "En bodega", Origen: o.Origen, Destino: o.Destino, Nombre: o.Nombre}
}

//CreaRegistro en el que escribira el camion.
func CreaRegistro() *os.File {
	csvFile, err := os.Create("RegistroLogistica.csv")

	if err != nil {
		log.Fatalf("Fallo al crear csv file: %s", err)
	}
	//Escribe lo que ira en cada columna
	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()
	val := []string{"timestamp", "id-paquete", "tipo", "nombre", "valor", "origen", "destino", "seguimiento"}
	csvwriter.Write(val)

	return csvFile

}

//EditaResigtro agrega registro del camion a el csv file.
func EditaResigtro(csvFile *os.File, o *sm.Orden, nSeg int32) {
	csvwriter := csv.NewWriter(csvFile)
	val := []string{time.Now().Format("2006-01-02 15:04:05"), o.Id, o.Tipo, o.Nombre, strconv.Itoa(int(o.Valor)), o.Origen, o.Destino, strconv.Itoa(int(nSeg))}
	csvwriter.Write(val)
}

//Para usar en local, cambiar ipport por ":"+port
func main() {
	// Escucha las conexiones grpc
	lis, err := net.Listen("tcp", ipportgrpc)

	if err != nil {
		log.Fatalf("Failed to listen on "+ipportgrpc+": %v", err)
	}

	s := Server{"1", []*sm.Paquete{}, []*sm.Paquete{}, []*sm.Paquete{}, make(map[int32]string)}

	//Crea el archivo de registro de logistica
	CreaRegistro()

	//Crea la conexion RabbitMQ
	conn, err := amqp.Dial(ipportrabbitmq)

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Inicializa el codigo de seguimiento
	CodSeg = 10000

	grpcServer := grpc.NewServer()

	fmt.Println("En espera de Informacion paquetes para servidor")

	//Inicia el servicio de mensajeria que contiene las funciones grpc
	sm.RegisterMensajeriaServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over "+ipportgrpc+": %v", err)
	}

}
