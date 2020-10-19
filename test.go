package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex
var wg sync.WaitGroup

func prints() {
	wg.Done()
	var t int
	var m map[int]string
	m = make(map[int]string)
	m[12] = "1"
	fmt.Println(m[12])
	fmt.Println(m[244] == "")
	fmt.Println("Ingrese tiempo de espera entre órdenes en segundos: ")
	fmt.Scanln(&t)
}

func main() {
	wg.Add(1)
	go prints()
	time.Sleep(2 * time.Second)
	wg.Add(1)
	go prints()
	wg.Wait()

}
