package main

import (
	"log"
	"math/rand"
	"reflect"

	"sync"

	sm "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup
var wg2 sync.WaitGroup

//IP local 10.6.40.161
const (
	//ipport = "10.6.40.162:50051"
	ipport = ":50051"
)

//Paquete Estructura del paquete a recibir.Tipo: retail, normal, prioritario. Estado: En bodega, en camino, recibido, no recibido.

//Camion Estructura de camion, se tienen 3 camiones. Tipo: retail, normal.
type Camion struct {
	id             int32
	tipo           string
	disponible     bool
	paq1           *sm.Paquete
	paq2           *sm.Paquete
	EntrPrevRetail bool
	PaqCargRetail  bool
}

//CamionResp para reconocer llamadas de camion en Logistica.
type CamionResp struct {
	id   int32
	tipo string
}

//EntregaPaquete intenta entregar paquete.
func EntregaPaquete() int {
	c := rand.Float64()
	if c < 0.8 {
		return 1
	}
	return 0

}

//ReintentaEntregar si es retail y ha intentado menos de 3 se puede reintentar, si es pyme depende del coste del producto y es a lo mas 2 reintentos.
func ReintentaEntregar(paq *sm.Paquete) int {
	if paq.Tipo == "Retail" {
		if paq.Intentos < 3 {
			return 1
		}
		return 0

	}
	if paq.Intentos < 2 {

		ganancia := float32(paq.Valor)
		if paq.Tipo == "Prioritario" {
			ganancia = ganancia + ganancia*0.3
		}
		ganancia = ganancia - float32(paq.Intentos+1)*10.0
		if ganancia > 0 {
			return 1
		}
		return 0

	}
	return 0

}

//IntentaEntregar retorna 0 si todas las entragas fueron exitosas, 1 si ninguna, 2 si solo la segunda y 3 si solo la primera.
func IntentaEntregar(paq *sm.Paquete, conn *grpc.ClientConn, ready bool) int {
	res := EntregaPaquete()
	if res == 1 && ready != true {
		paq.Estado = "Recibido"
		ready = true
	}
	//SendEstado
	wg.Add(1)
	go EntregaPosicionEntregaActual(conn, paq)
	return res
}

//CamionEntregaPaquetes intenta entregar los paquetes que lleva a sus destinos y vuelve a central.
func CamionEntregaPaquetes(cam *Camion, conn *grpc.ClientConn) {
	ready := false
	ready2 := false
	r1 := 0
	r2 := 0
	for ready != true && ready2 != true {
		if cam.paq1.Valor > cam.paq2.Valor {
			r1 = IntentaEntregar(cam.paq1, conn, ready)
			r2 = IntentaEntregar(cam.paq2, conn, ready2)
			if r1 == 1 && ready != true {
				cam.paq1.Estado = "Recibido"
				ready = true
			}
			if r2 == 1 && ready2 != true {
				cam.paq2.Estado = "Recibido"
				ready2 = true
			}

			if cam.paq1.Estado != "Recibido" {
				if ReintentaEntregar(cam.paq1) == 0 {
					ready = true
				}
			}
			if cam.paq2.Estado != "Recibido" {
				if ReintentaEntregar(cam.paq2) == 0 {
					ready2 = true
				}
			}
		} else {

		}
	}
}

//IniciaCliente inicia conexion cliente
func IniciaCliente() *grpc.ClientConn {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	return conn
}

//InformaPaqueteLogistica Camion informa entrega de la orden a Logistica
func InformaPaqueteLogistica(conn *grpc.ClientConn, cam Camion) string {
	defer wg.Done()
	c := sm.NewMensajeriaServiceClient(conn)
	ctx := context.Background()
	paquetes := []*sm.Paquete{}
	paquetes = append(paquetes, cam.paq1)

	if !reflect.DeepEqual(cam.paq2, sm.Paquete{}) {
		paquetes = append(paquetes, cam.paq2)
	}
	response, err := c.InformaEntrega(ctx, &sm.InformePaquetes{Paquetes: paquetes})
	if err != nil {
		log.Fatalf("Error al llamar InformaPaquete: %s", err)
	}

	log.Printf("Logistica: %s", response.Body) //Revisar mensaje de respuesta logistica.
	return response.Body
}

//EntregaPosicionEntregaActual Camion informa estado del paquete a Logistica
func EntregaPosicionEntregaActual(conn *grpc.ClientConn, paq *sm.Paquete) string {
	defer wg.Done()
	c := sm.NewMensajeriaServiceClient(conn)
	ctx := context.Background()

	response, err := c.EntregaPosicion(ctx, &sm.InformacionPaquete{Idpaquete: paq.Id, Estado: paq.Estado})
	if err != nil {
		log.Fatalf("Error al llamar EntregaPosicion: %s", err)
	}

	log.Printf("Logistica: %s", response.Body) //Revisar mensaje de respuesta logistica.
	return response.Body
}

//CamionDisponible Informa que el camion se encuentra disponible para cargar paquetes, ya sea tiene o no 1 paquete.
func CamionDisponible(conn *grpc.ClientConn, cam Camion) *sm.Paquete {
	defer wg.Done()
	c := sm.NewMensajeriaServiceClient(conn)
	ctx := context.Background()
	response, err := c.RecibeInstrucciones(ctx, &sm.DisponibleCamion{Id: cam.id, Tipo: cam.tipo, EntrPrevRetail: cam.EntrPrevRetail, PaqCargRetail: cam.PaqCargRetail})
	if err != nil {
		log.Fatalf("Error al llamar EntregaPosicion: %s", err)
	}
	return response

}

func holis(cam *Camion) {
	defer wg.Done()
	log.Printf("soy el camions %d\n", cam.id)

}

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	c1 := Camion{1, "Retail", true, &sm.Paquete{}, &sm.Paquete{}, false, false}
	c2 := Camion{2, "Retail", true, &sm.Paquete{}, &sm.Paquete{}, false, false}
	c3 := Camion{3, "Normal", true, &sm.Paquete{}, &sm.Paquete{}, false, false}

	wg.Add(1)
	go holis(&c1)
	wg.Add(1)
	go holis(&c2)
	wg.Add(1)
	go holis(&c3)
	wg.Add(1)
	go InformaPaqueteLogistica(conn, c1)
	wg.Add(1)
	go InformaPaqueteLogistica(conn, c2)
	wg.Add(1)
	go InformaPaqueteLogistica(conn, c3)
	wg.Wait()

}
