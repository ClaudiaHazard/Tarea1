package main

import (
	//"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
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

//IP local 10.6.40.163
const (
	ipport = "10.6.40.162:50051"
	//ipport = ":50051"
)

//EnviaOrdenCliente de Cliente a Logistica
func EnviaOrdenCliente(conn *grpc.ClientConn, tip string, aidi string, pro string, val int32, tien string, dest string) int32 {

	c := sm.NewMensajeriaServiceClient(conn)

	response, err2 := c.RealizaOrden(context.Background(), &sm.Orden{Id: aidi, Tipo: tip, Valor: val, Origen: tien, Destino: dest, Nombre: pro})

	if err2 != nil {
		log.Fatalf("Error al llamar RealizaOrden: %s", err2)
	}
	//log.Println("Orden registrada")

	return response.CodigoSeguimiento
}

//EnviaCodCliente cliente envia codigo de seguimiento
func EnviaCodCliente(conn *grpc.ClientConn, cod int32) string {

	c := sm.NewMensajeriaServiceClient(conn)

	response, err2 := c.SolicitaSeguimiento(context.Background(), &sm.CodSeguimiento{CodigoSeguimiento: cod})

	if err2 != nil {
		log.Fatalf("Error al llamar EnviaOrden: %s", err2)
	}
	//log.Printf("Código de sguimiento: %s", response.Body)

	return response.Estado
}

//IndividualOrder recibe una entrada del archic csv y la envía a logística
func IndividualOrder(record []string, tipo string, c *grpc.ClientConn) int32 {
	var order [6]string
	order[1] = record[0]
	order[2] = record[1]
	order[3] = record[2]
	order[4] = record[3]
	order[5] = record[4]
	if tipo == "retail" {
		order[0] = "retail"
	} else {
		if record[5] == "0" {
			order[0] = "normal"
		} else {
			order[0] = "prioritario"
		}
	}

	co, err := strconv.ParseInt(order[3], 10, 32)
	if err == nil {
		fmt.Println(co)
	}
	ccc := int32(co)
	rett := EnviaOrdenCliente(c, order[0], order[1], order[2], ccc, order[4], order[5])

	return rett

}

//Ordenar realiza orden de un cliente.
func Ordenar(tii string, c *grpc.ClientConn, pym [][]string, reta [][]string) {

	var ins []string
	if tii == "retail" {
		ins = reta[rand.Intn(len(reta)-2)+1]
		IndividualOrder(ins, tii, c)
		fmt.Println("Orden retail ingresada")
	} else {
		ins = pym[rand.Intn(len(pym)-2)+1]
		r := IndividualOrder(ins, tii, c)
		fmt.Printf("Orden pyme ingresada, este es su codigo de seguimiento: %d\n", r)
	}
}

//DoOrder recibe todos los datos de los csv, seleccionando uno al azar dependiendo del tipo de cliente
func DoOrder(pym [][]string, reta [][]string, c *grpc.ClientConn, m int) {
	defer wg.Done()
	var tii string
	for {
		locks="si"
		fmt.Println("Ingrese tipo de cliente (retail o pyme): ")
		fmt.Scanln(&tii)
		wg.Add(1)
		go Ordenar(tii, c, pym, reta)
		locks="no"
		time.Sleep(time.Duration(m) * time.Second)
	}
}

//PideSegui solicita ifnormación de un pquete con su código de seguimiento
func PideSegui(c *grpc.ClientConn) {
	defer wg.Done()
	time.Sleep(5 * time.Second)
	for {
		if locks!="si"{
			var codd int32
			fmt.Println("Ingrese codigo de seguimiento: ")
			fmt.Scanln(&codd)

			//envío y recepción de info de estado
			info := EnviaCodCliente(c, codd)
			//mostrar info
			fmt.Println("Estado del paquete: ", info)
		}
	}
}
func main() {

	var conn *grpc.ClientConn
	var t int
	var fx string
	var fx2 string

	conn, err2 := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())
	locks ="si"

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
		log.Fatalln("Couldn't open the csv file", err)
	}
	csvfile2, err2 := os.Open(fx2)
	if err2 != nil {
		log.Fatalln("Couldn't open the csv file", err)
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

	wg.Add(1)
	go DoOrder(allpyme, allretail, conn, t)
	wg.Add(1)
	go PideSegui(conn)

	wg.Wait()
}
