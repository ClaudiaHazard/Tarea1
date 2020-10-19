package main

import (
	"fmt"
	"math/rand"
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

func ranss() {
	c := rand.Float64()
	if c < 0.8 {
		fmt.Println(1)
	} else {
		fmt.Println(0)
	}
}

func main() {
	wg.Add(1)
	go prints()
	time.Sleep(2 * time.Second)
	wg.Add(1)
	go prints()
	wg.Wait()

	for range []int{6, 2, 4, 4, 45} {
		ranss()
	}

}
