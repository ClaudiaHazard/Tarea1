package main
import (
	//"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
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
    fx=fx + ".csv" 
    t=1
	csvfile, err := os.Open(fx)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	t=1
	r := csv.NewReader(csvfile)
    if cli=="retail"  {
    	var a int
    	a=0
		for {
			// Read each record from csv

			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if a!=0 {
				order[0]="retail"
				order[1]=record[0]
				order[2]=record[1]
				order[3]=record[2]
				order[4]=record[3]
			//comunicarla al logistica y RECIBIR COD DE VERIFICACIÓN
			}
			//sleep
			time.Sleep(time.Duration(t) * time.Second) 
			fmt.Println(order) 
			a=a+1

		}
    } else {
    	//otro tipo
    	var a int
    	a=0
    	for {
			// Read each record from csv
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if a!=0{
					order[1]=record[0]
					order[2]=record[1]
					order[3]=record[2]
					order[4]=record[3]
				if record[4]=="0"{
					order[0]="normal"
				}else {
					order[0]="prioritario"
				}
			}
			//comunicarla al logistica y RECIBIR COD DE VERIFICACIÓN

			//sleep
			time.Sleep(time.Duration(t) * time.Second) 
			fmt.Println(order) 
			a=a+1
			
		}
    }
    //Seguimiento de órdenes
    for{
	    var cod string
	    fmt.Println("Ingrese codigo de seguimiento: ") 
	    fmt.Scanln(&cod)

	    //envío y recepción de info de estado
	}
}