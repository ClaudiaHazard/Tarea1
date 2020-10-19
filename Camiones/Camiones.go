package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

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

//Camion Estructura de camion, se tienen 3 camiones. Tipo: retail, normal.
type Camion struct {
	id             int32
	tipo           string
	disponible     bool
	paq1           *sm.Paquete
	fechaEntrega1  string
	paq2           *sm.Paquete
	fechaEntrega2  string
	EntrPrevRetail bool
	PaqCargRetail  bool
}

//CamionResp para reconocer llamadas de camion en Logistica.
type CamionResp struct {
	id   int32
	tipo string
}

//EntregaPaquete intenta entregar paquete.
func EntregaPaquete(te int) int {
	time.Sleep(time.Duration(te) * time.Millisecond)
	c := rand.Float64()
	if c < 0.8 {
		return 1
	}
	return 0

}

//ReintentaEntregar si es retail y ha intentado menos de 3 se puede reintentar, si es pyme depende del coste del producto y es a lo mas 2 reintentos.
func ReintentaEntregar(paq *sm.Paquete) int {
	if paq.Tipo == "retail" {
		if paq.Intentos < 3 {
			return 1
		}
		return 0

	}
	if paq.Intentos < 2 {

		ganancia := float32(paq.Valor)
		if paq.Tipo == "prioritario" {
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
func IntentaEntregar(paq *sm.Paquete, conn *grpc.ClientConn, ready bool, te int) (int, string) {

	res := EntregaPaquete(te)
	tiempoEntrega := "0"
	if res == 1 && ready != true {
		paq.Estado = "Recibido"
		ready = true
		tiempoEntrega = time.Now().Format("2006-01-02 15:04:05")

	}
	//SendEstado
	go EntregaPosicionEntregaActual(conn, paq)
	return res, tiempoEntrega
}

//CamionEntregaPaquetes intenta entregar los paquetes que lleva a sus destinos y vuelve a central.
func CamionEntregaPaquetes(cam *Camion, conn *grpc.ClientConn, te int) {
	ready := false
	ready2 := false
	r1 := 0
	r2 := 0
	t1 := "0"
	t2 := "0"
	for ready != true && ready2 != true {
		//Entrega el mas caro primero
		if cam.paq1.Valor > cam.paq2.Valor {
			r1, t1 = IntentaEntregar(cam.paq1, conn, ready, te)

			r2, t2 = IntentaEntregar(cam.paq2, conn, ready2, te)

			if r1 == 1 && ready != true {
				cam.paq1.Estado = "Recibido"
				cam.fechaEntrega1 = t1
				ready = true

			}
			if r2 == 1 && ready2 != true {
				cam.paq2.Estado = "Recibido"
				cam.fechaEntrega2 = t2
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
			r1, t1 = IntentaEntregar(cam.paq2, conn, ready, te)
			r2, t2 = IntentaEntregar(cam.paq1, conn, ready2, te)
			if r1 == 1 && ready != true {
				cam.paq2.Estado = "Recibido"
				cam.fechaEntrega2 = t1
				ready = true
			}
			if r2 == 1 && ready2 != true {
				cam.paq1.Estado = "Recibido"
				cam.fechaEntrega1 = t2
				ready2 = true
			}

			if cam.paq2.Estado != "Recibido" {
				if ReintentaEntregar(cam.paq2) == 0 {
					ready = true
				}
			}
			if cam.paq1.Estado != "Recibido" {
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

//InformaPaqueteLogistica Camion informa entrega de la orden a Logistica
func InformaPaqueteLogistica(conn *grpc.ClientConn, cam *Camion) string {
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
	c := sm.NewMensajeriaServiceClient(conn)
	ctx := context.Background()

	response, err := c.EntregaPosicion(ctx, &sm.InformacionPaquete{CodigoSeguimiento: paq.CodigoSeguimiento, Estado: paq.Estado})
	if err != nil {
		log.Fatalf("Error al llamar EntregaPosicion: %s", err)
	}

	log.Printf("Logistica: %s", response.Body) //Revisar mensaje de respuesta logistica.
	return response.Body
}

//CamionDisponible Informa que el camion se encuentra disponible para cargar paquetes, ya sea tiene o no 1 paquete.
func CamionDisponible(conn *grpc.ClientConn, cam *Camion) *sm.Paquete {
	c := sm.NewMensajeriaServiceClient(conn)
	ctx := context.Background()
	response, err := c.RecibeInstrucciones(ctx, &sm.DisponibleCamion{Id: cam.id, Tipo: cam.tipo, EntrPrevRetail: cam.EntrPrevRetail, PaqCargRetail: cam.PaqCargRetail})
	if !cam.PaqCargRetail && cam.tipo == "Retail" && cam.PaqCargRetail == false {
		cam.PaqCargRetail = true
	}
	if err != nil {
		log.Fatalf("Error al llamar EntregaPosicion: %s", err)
	}
	return response
}

//ComparaPaquete  para que sepa si es vacio porque las funciones de comparacion no sirven.
func ComparaPaquete(paq1 *sm.Paquete, paq2 *sm.Paquete) bool {
	if paq1.CodigoSeguimiento != paq2.CodigoSeguimiento {
		return false
	}
	if paq1.Destino != paq2.Destino {
		return false
	}
	if paq1.Estado != paq2.Estado {
		return false
	}
	if paq1.Id != paq2.Id {
		return false
	}
	if paq1.Intentos != paq2.Intentos {
		return false
	}
	if paq1.Nombre != paq2.Nombre {
		return false
	}
	if paq1.Origen != paq2.Origen {
		return false
	}
	if paq1.Destino != paq2.Destino {
		return false
	}
	if paq1.Tipo != paq2.Tipo {
		return false
	}
	if paq1.Valor != paq2.Valor {
		return false
	}
	return true
}

//CamionEspera camión que no tiene paquetes recibe un paquete, y luego espera a poder cargar uno
func CamionEspera(cam *Camion, conn *grpc.ClientConn, ti int, te int) {
	csv := CreaRegistro(cam)
	EmptyPaq := &sm.Paquete{}
	defer wg.Done()
	tRec := time.Now()
	//Por ahora ciclo infinito simplemente
	for {
		log.Printf("Camion %d en espera de paquete", cam.id)
		for ComparaPaquete(cam.paq1, EmptyPaq) {
			cam.paq1 = CamionDisponible(conn, cam)
			tRec = time.Now().Add(time.Millisecond * time.Duration(ti))
		}
		log.Printf("Camion %d recibe paquete 1, espera por paquete 2.", cam.id)
		for ComparaPaquete(cam.paq2, EmptyPaq) && (tRec.Sub(time.Now()) > time.Duration(0)) {
			cam.paq2 = CamionDisponible(conn, cam)
			//Por ahora lo deje así para que no mande tantos mensajes.
			time.Sleep(500 * time.Millisecond)
		}
		cam.disponible = false
		//Paquetes salen de central
		cam.paq1.Estado = "En Camino"
		log.Printf("Camion %d sale de central con paquete 1 con id %s.", cam.id, cam.paq1.Id)
		if !ComparaPaquete(cam.paq2, &sm.Paquete{}) {
			cam.paq1.Estado = "En Camino"
			log.Printf("Camion %d sale de central con paquete 2 con id %s.", cam.id, cam.paq2.Id)
		}

		CamionEntregaPaquetes(cam, conn, te)
		//Si vuelve a central y los paquetes no cambiaron su estado a recibido, el paquete no fue entregado.
		if cam.paq1.Estado == "En Camino" {
			cam.paq1.Estado = "No Recibido"
		}
		if cam.paq2.Estado == "En Camino" {
			cam.paq2.Estado = "No recibido"
		}

		//Camion avisa que vuelve a central.
		log.Printf("Camion %d vuelve a central.", cam.id)
		InformaPaqueteLogistica(conn, cam)

		EditaResigtro(cam, csv)
		cam = VaciaCamion(cam)
		//Camion vuelve a avisar que esta en espera de mas paquetes.
	}

}

//VaciaCamion vacia los paquetes en central y vuelve a estar disponible
func VaciaCamion(cam *Camion) *Camion {
	cam.paq1 = &sm.Paquete{}
	cam.paq2 = &sm.Paquete{}
	cam.disponible = true
	cam.fechaEntrega1 = "0"
	cam.fechaEntrega2 = "0"
	return cam
}

//CreaRegistro en el que escribira el camion.
func CreaRegistro(cam *Camion) *os.File {
	csvFile, err := os.Create("RegistroCamion" + strconv.Itoa(int(cam.id)) + ".csv")

	if err != nil {
		log.Fatalf("Fallo al crear csv file: %s", err)
	}

	//Escribe lo que ira en cada columna
	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()
	val := []string{"id-paquete", "tipo", "valor", "origen", "destino", "intentos", "fecha-entrega"}
	csvwriter.Write(val)

	return csvFile

}

//EditaResigtro agrega registro del camion a el csv file.
func EditaResigtro(cam *Camion, csvFile *os.File) {
	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()
	val := []string{cam.paq1.Id, cam.paq1.Tipo, strconv.Itoa(int(cam.paq1.Valor)), cam.paq1.Origen, cam.paq1.Destino, strconv.Itoa(int(cam.paq1.Intentos)), cam.fechaEntrega1}
	csvwriter.Write(val)
	if (cam.paq2 != &sm.Paquete{}) {
		val = []string{cam.paq2.Id, cam.paq2.Tipo, strconv.Itoa(int(cam.paq2.Valor)), cam.paq2.Origen, cam.paq2.Destino, strconv.Itoa(int(cam.paq2.Intentos)), cam.fechaEntrega2}
		csvwriter.Write(val)
	}
}

func main() {
	var conn *grpc.ClientConn
	var ti int
	var te int

	fmt.Println("Ingrese duración de espera por segunda orden en milisegundos: ")
	fmt.Scanln(&ti)

	fmt.Println("Ingrese duración de entrega por paquete en milisegundos: ")
	fmt.Scanln(&te)

	//Se crea la conexion con el servidor Logistica
	conn, err := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("No se pudo conectar: %s", err)
	}

	//Se inicializan los 3 camiones.
	c1 := Camion{1, "Retail", true, &sm.Paquete{}, "0", &sm.Paquete{}, "0", false, false}
	c2 := Camion{2, "Retail", true, &sm.Paquete{}, "0", &sm.Paquete{}, "0", false, false}
	c3 := Camion{3, "Normal", true, &sm.Paquete{}, "0", &sm.Paquete{}, "0", false, false}

	wg.Add(1)
	go CamionEspera(&c1, conn, ti, te)
	wg.Add(1)
	go CamionEspera(&c2, conn, ti, te)
	wg.Add(1)
	go CamionEspera(&c3, conn, ti, te)

	wg.Wait()

}
