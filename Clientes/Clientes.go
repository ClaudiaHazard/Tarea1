package main

import (
	//"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	sm "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup
var mutex sync.Mutex
var locks string

//AunExistenPaquetes Variable para determinar si quedan paquetes en los archivos csv.
var AunExistenPaquetes bool

//CountRetail variable que contiene la cantidad de paquetes leidos en el csv.
var CountRetail int

//CountPyme variable que contiene la cantidad de paquetes leidos en el csv.
var CountPyme int

//IP local 10.6.40.163
const (
	//ipport = "10.6.40.162:50051"
	ipport = ":50051"
)

//EnviaOrdenCliente de Cliente a Logistica
func EnviaOrdenCliente(conn *grpc.ClientConn, ord *sm.Orden) {
	c := sm.NewMensajeriaServiceClient(conn)

	response, err2 := c.RealizaOrden(context.Background(), ord)

	if err2 != nil {
		log.Fatalf("Error al llamar RealizaOrden: %s", err2)
	}
	if response.CodigoSeguimiento == 0 {
		log.Printf("Orden creada, usted no cuenta con numero de seguimiento\n")
	} else {
		log.Printf("Orden creada, su numero de seguimiento es: %d\n", response.CodigoSeguimiento)
	}

}

//EnviaCodCliente cliente envia codigo de seguimiento
func EnviaCodCliente(conn *grpc.ClientConn, cod int32) string {

	c := sm.NewMensajeriaServiceClient(conn)

	response, err2 := c.SolicitaSeguimiento(context.Background(), &sm.CodSeguimiento{CodigoSeguimiento: cod})

	if err2 != nil {
		log.Fatalf("Error al llamar EnviaOrden: %s", err2)
	}
	if response.Estado == "" {
		log.Printf("Su codigo %d no corresponde a ningun paquete enviado.", cod)
	} else {
		log.Printf("Su paquete con codigo %d, se encuentra en estado: %s", cod, response.Estado)
	}

	return response.Estado
}

//GeneraOrdenFromString genera la orden del string del csv
func GeneraOrdenFromString(record []string, tipo string) *sm.Orden {
	var order [6]string
	order[1] = record[0]
	order[2] = record[1]
	order[3] = record[2]
	order[4] = record[3]
	order[5] = record [4]
	if tipo == "retail" {
		order[0] = "retail"
	} else {
		if record[5] == "0" {
			order[0] = "normal"
		} else {
			order[0] = "prioritario"
		}
	}

	val, _ := strconv.ParseInt(order[3], 10, 32)

	return &sm.Orden{Id: order[1], Tipo: order[0], Valor: int32(val), Origen: order[4], Destino: order[5], Nombre: order[2]}

}

func getValue(csv [][]string, count int) ([]string, int) {
	defer mutex.Unlock()
	mutex.Lock()
	val := csv[count]
	count = count + 1
	return val, count
}

//CreaOrdenTipo Crea orden del documento Pyme o Retail
func CreaOrdenTipo(tipo string, conn *grpc.ClientConn, CsvPyme [][]string, CsvRetail [][]string) {
	var ordCsv []string
	if tipo == "retail" {
		if len(CsvRetail) > CountRetail {
			ordCsv, CountRetail = getValue(CsvRetail, CountRetail)
			ord := GeneraOrdenFromString(ordCsv, "retail")
			EnviaOrdenCliente(conn, ord)
			log.Printf("ContadorOrdenes Retail: %d\n", CountRetail)

		} else {
			fmt.Println("No quedan mas paquetes en el archivo csv de Retail")
		}

	}
	if tipo == "pyme" {
		if len(CsvPyme) > CountPyme {
			mutex.Lock()
			ordCsv = CsvPyme[CountPyme]
			ord := GeneraOrdenFromString(ordCsv, "pyme")
			CountPyme = CountPyme + 1
			mutex.Unlock()
			EnviaOrdenCliente(conn, ord)
			log.Printf("ContadorOrdenes Pyme: %d\n", CountPyme)
		} else {
			fmt.Println("No quedan mas paquetes en el archivo csv de Pymes")
		}
	}
}

func main() {

	var conn *grpc.ClientConn
	var t int
	var fx string
	var fx2 string

	conn, err2 := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())
	locks = "si"

	if err2 != nil {
		log.Fatalf("did not connect: %s", err2)
	}
	defer conn.Close()

	fmt.Println("Ingrese tiempo de espera entre órdenes en segundos: ")

	fmt.Scanln(&t)
	fmt.Println("Ingrese nombre de archivo csv retail(ejemplo: si es retail.csv usted escribe retail): ")
	fmt.Scanln(&fx)
	fx = fx + ".csv"
	fmt.Println("Ingrese nombre de archivo csv pymes(ejemplo: si es pymes.csv usted escribe pymes): ")
	fmt.Scanln(&fx2)
	fx2 = fx2 + ".csv"
	csvfile, err := os.Open(fx)
	if err != nil {
		log.Fatalln("No se pudo abrir el archivo csv ", err)
	}
	csvfile2, err2 := os.Open(fx2)
	if err2 != nil {
		log.Fatalln("No se pudo abrir el archivo csv", err)
	}
	r := csv.NewReader(csvfile)
	r2 := csv.NewReader(csvfile2)

	allretail, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	allpyme, err := r2.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	CountPyme = 1
	CountRetail = 1
	AunExistenPaquetes = true

	var tipo string
	var cod int32
	tOrden := time.Now()

	for AunExistenPaquetes {
		for tipo != "retail" && tipo != "pyme" {
			time.Sleep(1 * time.Second)
			fmt.Println("Ingrese tipo de cliente (retail o pyme): ")
			fmt.Scanln(&tipo)
			if tipo != "retail" && tipo != "pyme" {
				fmt.Println("Valor erroneo")
			} else {
				go CreaOrdenTipo(tipo, conn, allpyme, allretail)
				tOrden = time.Now().Add(time.Millisecond * time.Duration(t))
			}

		}
		for tOrden.Sub(time.Now()) > time.Duration(0) {
			time.Sleep(1 * time.Second)
			fmt.Println("Ingrese codigo de seguimiento: ")
			fmt.Scanln(&cod)
			if cod != 0 {
				go EnviaCodCliente(conn, cod)
			} else {
				fmt.Println("No existe registro de ese codigo de seguimiento")
			}
		}
		tipo = ""
	}
}
