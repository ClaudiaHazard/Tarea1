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

	enviaorden "github.com/ClaudiaHazard/Tarea1/Logistica/EnviaOrden"
	"google.golang.org/grpc"
)

//IP local 10.6.40.163
const (
	ipport = "10.6.40.162:50052"
	//ipport= ":50052"
)

func main() {

	fmt.Println("Ingrese tipo de cliente: ")
	var cli string
	var fx string
	var order [5]string
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
				//comunicarla al logistica y RECIBIR COD DE VERIFICACIÓN
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
				if record[4] == "0" {
					order[0] = "normal"
				} else {
					order[0] = "prioritario"
				}
			}
			//comunicarla al logistica y RECIBIR COD DE VERIFICACIÓN

			//sleep
			time.Sleep(time.Duration(t) * time.Second)
			fmt.Println(order)
			a = a + 1

		}
	}
	//Seguimiento de órdenes
	for {
		var cod string
		fmt.Println("Ingrese codigo de seguimiento: ")
		fmt.Scanln(&cod)

		//envío y recepción de info de estado
	}

	fmt.Println("Conexion con Logistica")

	var conn *grpc.ClientConn

	conn, err2 := grpc.Dial(ipport, grpc.WithInsecure(), grpc.WithBlock())

	if err2 != nil {
		log.Fatalf("did not connect: %s", err2)
	}
	defer conn.Close()

	fmt.Println("Crea conexion para enviar orden")

	c := enviaorden.NewEnviaOrdenServiceClient(conn)

	fmt.Println("Envia Mensaje")

	response, err2 := c.SayHello(context.Background(), &enviaorden.Message{Body: "Hola por parte de Cliente!"})

	if err2 != nil {
		log.Fatalf("Error al llamar SayHello: %s", err2)
	}
	log.Printf("Respuesta de Logistica: %s", response.Body)

}
