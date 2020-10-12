package main

import (
	servidorcamiones "github.com/ClaudiaHazard/Tarea1/Logistica/ServidorCamiones"
	servidorcliente "github.com/ClaudiaHazard/Tarea1/Logistica/ServidorCliente"
)

//SayHello Test
func main() {
	servidorcliente.IniciarServidorCliente()
	servidorcamiones.IniciarServidorCamiones()
}
