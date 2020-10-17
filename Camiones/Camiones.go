package main

import (
	"log"
	"math/rand"

	"sync"

	sm "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

//IP local 10.6.40.161
const (
	//ipport = "10.6.40.162:50051"
	ipport = ":50051"
)

//Paquete Estructura del paquete a recibir.Tipo: retail, normal, prioritario. Estado: En bodega, en camino, recibido, no recibido.
type Paquete struct {
	id                int32
	CodigoSeguimiento int32
	tipo              string
	valor             int32
	intentos          int32
	estado            string
	origen            string
	destino           string
}

//Camion Estructura de camion, se tienen 3 camiones. Tipo: retail, normal.
type Camion struct {
	id         int32
	tipo       string
	disponible bool
	paq1       Paquete
	paq2       Paquete
}

//EntregaPaquete intenta entregar paquete.
func EntregaPaquete(paq Paquete) int {
	c := rand.Float64()
	if c < 0.8 {
		return 1
	} else {
		return 0
	}
}

//ReintentaEntregar si es retail y ha intentado menos de 3 se puede reintentar, si es pyme depende del coste del producto y es a lo mas 2 reintentos.
func ReintentaEntregar(paq Paquete) int {
	if paq.tipo == "Retail" {
		if paq.intentos < 3 {
			return 1
		} else {
			return 0
		}
	} else {
		if paq.intentos < 2 {

			ganancia := float32(paq.valor)
			if paq.tipo == "Prioritario" {
				ganancia = ganancia + ganancia*0.3
			}
			ganancia = ganancia - float32(paq.intentos+1)*10.0
			if ganancia > 0 {
				return 1
			} else {
				return 0
			}

		} else {
			return 0
		}
	}

}

//IntentaEntregar retorna 0 si todas las entragas fueron exitosas, 1 si ninguna, 2 si solo la segunda y 3 si solo la primera.
func IntentaEntregar(paq1 Paquete, paq2 Paquete) (int, int) {
	res := EntregaPaquete(paq1)
	res2 := EntregaPaquete(paq2)
	return res, res2
}

//CamionEntregaPaquetes intenta entregar los paquetes que lleva a sus destinos y vuelve a central.
func CamionEntregaPaquetes(cam *Camion) {
	ready := false
	ready2 := false
	for ready != true && ready2 != true {
		if cam.paq1.valor > cam.paq2.valor {
			r1, r2 := IntentaEntregar(cam.paq1, cam.paq2)
			if r1 == 1 && ready != true {
				cam.paq1.estado = "Recibido"
				ready = true
			}
			if r2 == 1 && ready2 != true {
				cam.paq2.estado = "Recibido"
				ready2 = true
			}
			//SendEstado

			if cam.paq1.estado != "Recibido" {
				if ReintentaEntregar(cam.paq1) == 0 {
					ready = true
				}
			}
			if cam.paq2.estado != "Recibido" {
				if ReintentaEntregar(cam.paq1) == 0 {
					ready2 = true
				}
			}
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

//InformaPaqueteLogistica Camion informa estado del paquete a Logistica
func InformaPaqueteLogistica(conn *grpc.ClientConn) string {
	c := sm.NewMensajeriaServiceClient(conn)
	ctx := context.Background()
	response, err := c.InformaEntrega(ctx, &sm.Message{Body: "Hola por parte de Camiones!"})
	if err != nil {
		log.Fatalf("Error al llamar InformaPaquete: %s", err)
	}

	log.Printf("Respuesta de Logistica: %s", response.Body)
	return response.Body
}

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	c1 := Camion{1, "Retail", true, Paquete{}, Paquete{}}
	c2 := Camion{2, "Retail", true, Paquete{}, Paquete{}}
	c3 := Camion{3, "Normal", true, Paquete{}, Paquete{}}

	go InformaPaqueteLogistica(conn)
	go InformaPaqueteLogistica(conn)
	go InformaPaqueteLogistica(conn)

}
