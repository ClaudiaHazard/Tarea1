package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/streadway/amqp"
)

//IP local 10.6.40.164
const (
	ipportrabbitmq = "amqp://test:test@10.6.40.162:5672/"
	//ipportrabbitmq = "amqp://guest:guest@localhost:5672/"
)

//Entry Struct de entrada de Finanzas
type Entry struct {
	ID        string
	Intentos  int
	Entregado bool
	Valor     int
	Tipo      string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//SumTotal obtiene total de ganancias
func SumTotal(prii int, tries int, reci bool, ret string) float64 {
	var sum float64
	sum = sum - float64(tries*10)
	if reci {
		sum = sum + float64(prii)
	} else {
		if ret == "prioritario " {
			sum = sum + (0.3 * float64(prii))
		} else if ret == "retail" {
			sum = sum + float64(prii)
		}
	}
	return sum
}

func getLastLineWithSeek(filepath string) []string {
	csvfile, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	return records[len(records)-1]
}

//cambiar gusername y password en vm
func main() {
	conn, err := amqp.Dial(ipportrabbitmq)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close() //"fmt"

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//crear el csv si no existe

	if _, err := os.Stat("caja.csv"); os.IsNotExist(err) {
		file, err := os.Create("caja.csv")
		if err != nil {
			log.Fatalln("Couldn't open the csv file", err)
		}
		x := []string{"Numero", "ID", "Intentos", "Entregado", "Ingresos", "Total"}
		csvWriter := csv.NewWriter(file)
		strWrite := [][]string{x}
		csvWriter.WriteAll(strWrite)
		csvWriter.Flush()
		file.Close()
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			//acá debería comprobarse el tipo de msj. hay que saber si se terminó o no el programa

			mess := d.Body

			var entry Entry
			json.Unmarshal([]byte(mess), &entry)

			//luego de recibrlo, se debe procesar
			var count int
			var aidi string
			var tri int
			var recei bool
			//var prie int
			var tip string
			var tot float64
			var TT float64
			var unit int
			aidi = entry.ID
			tri = entry.Intentos
			recei = entry.Entregado
			tip = entry.Tipo
			unit = entry.Valor
			line := getLastLineWithSeek("caja.csv")
			if aidi != "0" && aidi != "" {

				//int prii, int tries, bool reci, string ret
				//cálculo de ganancia/pérdida
				tot = SumTotal(unit, tri, recei, tip)
				if line[0] != "Número" {
					count, err = strconv.Atoi(line[0])
					count = count + 1
					TT, err = strconv.ParseFloat(line[5], 64)
					TT = TT + tot
				} else {
					count = 1
					TT = tot
				}

				csvfile, err := os.OpenFile("caja.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
				if err != nil {
					log.Fatalln("Couldn't open the csv file", err)
				}
				x := []string{strconv.Itoa(count), aidi, strconv.Itoa(tri), strconv.FormatBool(recei), fmt.Sprintf("%f", tot), fmt.Sprintf("%f", TT)}
				csvWriter := csv.NewWriter(csvfile)
				strWrite := [][]string{x}
				//fmt.Println(strWrite)
				csvWriter.WriteAll(strWrite)
				csvWriter.Flush()
				csvfile.Close()
			} else if aidi=="0" {
				line := getLastLineWithSeek("caja.csv")
				fmt.Println("Total hasta ahora: ", line[5])
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
