package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	sm "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
)

var wg sync.WaitGroup
var wg2 sync.WaitGroup

//Server Struct que contiene los valores del server
type Server struct {
	camion         int
	arrRetail      []Paquete
	arrPrioritario []Paquete
	arrNormal      []Paquete
}

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

//AgregaACola agrega paquete a cola correspondiente
func AgregaACola(p Paquete, s Server) {
	if p.tipo == "Retail" {
		s.arrRetail = append(s.arrRetail, p)
	}
	if p.tipo == "Prioritario" {
		s.arrPrioritario = append(s.arrPrioritario, p)
	}
	if p.tipo == "Retail" {
		s.arrNormal = append(s.arrNormal, p)
	}
}

//BorrarElemento borra el elemento en la posicion pos.
func BorrarElemento(arr []Paquete, pos int) []Paquete {
	copy(arr[pos:], arr[pos+1:]) // Shift a[i+1:] left one index.
	arr[len(arr)-1] = Paquete{}  // Erase last element (write zero value).
	arr = arr[:len(arr)-1]
	return arr
}

//AsignaPaquete asigna paquete al tipo de camion correspondiente.
func AsignaPaquete(s *Server, tipoCam string, entrPrevRetail bool, paqCargRetail bool) Paquete {
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
			return Paquete{}
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
			return Paquete{}
		}
	}
	return Paquete{}
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

func holis(cam *Camion) {
	defer wg.Done()
	log.Printf("soy el camions %d\n", cam.id)

}

func holis2(cam *Camion) {
	defer wg.Done()
	log.Printf("soy el camions %d\n", cam.id)
	time.Sleep(2 * time.Second)
	log.Printf("soy el camions %d\n", cam.id)
}

func inicializa(cam *Camion) {
	defer wg.Done()
	wg.Add(1)
	go holis(cam)
	wg.Add(1)
	go holis2(cam)

}

//Borrarpos borra elemento en posicion pos
func Borrarpos(arr []*sm.Paquete, pos int) []*sm.Paquete {
	copy(arr[pos:], arr[pos+1:])    // Shift a[i+1:] left one index.
	arr[len(arr)-1] = &sm.Paquete{} // Erase last element (write zero value).
	arr = arr[:len(arr)-1]
	return arr
}

func main() {

	c1 := Camion{1, "Retail", true, Paquete{}, Paquete{}}
	c2 := Camion{2, "Retail", true, Paquete{}, Paquete{}}
	c3 := Camion{3, "Normal", true, Paquete{}, Paquete{}}

	wg.Add(1)
	go inicializa(&c1)
	log.Printf("Termino w2")

	wg.Add(1)
	go inicializa(&c2)
	wg.Add(1)
	go inicializa(&c3)

	m := make(map[int]string)
	m[123] = "En proceso"
	log.Printf("Estado del paquete con codigo %s\n", m[123])

	ti := 5
	tRec := time.Now().Add(time.Second * time.Duration(ti))
	log.Println(tRec)
	log.Println(tRec.Sub(time.Now()))
	for tRec.Sub(time.Now()) > time.Duration(0) {
		log.Printf("Hola\n")
		time.Sleep(500 * time.Millisecond)
	}

	paq1 := &sm.Paquete{Id: "1", CodigoSeguimiento: 1, Tipo: "Retail", Valor: 10, Intentos: 0, Estado: "En bodega", Origen: "Origen A", Destino: "Destino A", Nombre: "Bicicleta"}
	paq2 := &sm.Paquete{Id: "2", CodigoSeguimiento: 2, Tipo: "Retail", Valor: 10, Intentos: 0, Estado: "En bodega", Origen: "Origen A", Destino: "Destino A", Nombre: "Bicicleta"}

	arrPaq := []*sm.Paquete{}
	arrPaq = append(arrPaq, paq1)
	arrPaq = append(arrPaq, paq2)
	log.Println(arrPaq)
	Borrarpos(arrPaq, 0)
	log.Println(arrPaq)

	wg.Wait()
}
