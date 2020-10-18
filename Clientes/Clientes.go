package main

import (
	//"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"strconv"

	sm "github.com/ClaudiaHazard/Tarea1/ServicioMensajeria"
	"google.golang.org/grpc"
)

//IP local 10.6.40.163
const (
	//ipport = "10.6.40.162:50051"
	ipport = ":50051"
)

//EnviaOrdenCliente de Cliente a Logistica
func EnviaOrdenCliente(conn *grpc.ClientConn, tip string, aidi string, pro string, val int32, tien string, dest string) int32 {

	c := sm.NewMensajeriaServiceClient(conn)

	response, err2 := c.RealizaOrden(context.Background(), &sm.Orden{Id: aidi, Tipo: tip, Valor:val, Origen: tien, Destino: dest,Nombre:pro})

	if err2 != nil {
		log.Fatalf("Error al llamar EnviaOrden: %s", err2)
	}
	log.Println("Orden registrada")

	return response.CodigoSeguimiento
}

//EnviaOrdenCliente de Cliente a Logistica
func EnviaCodCliente(conn *grpc.ClientConn, cod int32) string {

	c := sm.NewMensajeriaServiceClient(conn)

	response, err2 := c.SolicitaSeguimiento(context.Background(), &sm.CodSeguimiento{CodigoSeguimiento: cod})

	if err2 != nil {
		log.Fatalf("Error al llamar EnviaOrden: %s", err2)
	}
	//log.Printf("Código de sguimiento: %s", response.Body)

	return response.Estado
}


func main() {

	var conn *grpc.ClientConn

	conn, err2 := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())

	if err2 != nil {
		log.Fatalf("did not connect: %s", err2)
	}
	defer conn.Close()

	fmt.Println("Ingrese tipo de cliente: ")
	var cli string
	var fx string
	var order [6]string
	var t int
	fmt.Scanln(&cli)
	fmt.Println("Ingrese nombre de archivo: ")
	fmt.Scanln(&fx)
	fx = fx + ".csv"
	t = 1
	csvfile, err := os.Open(fx)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	t = 1
	r := csv.NewReader(csvfile)
	if cli == "retail" {
		var a int
		a = 0
		for {
			// Read each record from csv

			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if a != 0 {
				order[0] = "retail"
				order[1] = record[0] 
				order[2] = record[1]
				order[3] = record[2]
				order[4] = record[3]
				order[5] = record[4]
				//comunicarla al logistica y RECIBIR COD DE VERIFICACIÓN
				co,err:=strconv.ParseInt(order[3],10,32)
				if err == nil {
					fmt.Println(co)
				}
				c:=int32(co)
				EnviaOrdenCliente(conn,order[0],order[1],order[2],c,order[4],order[5])
				//tipo,id,prod,valor,tienda,destino
			}
			//sleep
			time.Sleep(time.Duration(t) * time.Second)
			fmt.Println(order)
			a = a + 1

		}
	} else {
		//otro tipo
		var a int
		a = 0
		for {
			// Read each record from csv
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if a != 0 {
				order[1] = record[0]
				order[2] = record[1]
				order[3] = record[2]
				order[4] = record[3]
				order[5] = record[4]
				if record[5] == "0" {
					order[0] = "normal"
				} else {
					order[0] = "prioritario"
				}
				//comunicarla al logistica y RECIBIR COD DE VERIFICACIÓN
				co,err:=strconv.ParseInt(order[3],10,32)
				if err == nil {
					fmt.Println(co)
				}
				c:=int32(co)				
				rett := EnviaOrdenCliente(conn,order[0],order[1],order[2],c,order[4],order[5])
				fmt.Println("Su código de seguimiento es: ",rett)
			}
			//sleep
			time.Sleep(time.Duration(t) * time.Second)
			fmt.Println(order)
			a = a + 1

		}
	}
	//Seguimiento de órdenes
	for {
		var codd int32
		fmt.Println("Ingrese codigo de seguimiento: ")
		fmt.Scanln(&codd)

		//envío y recepción de info de estado
		info := EnviaCodCliente(conn,codd)
		//mostrar info
		fmt.Println("Estado del paquete: ",info)
	}

}
