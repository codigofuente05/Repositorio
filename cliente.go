package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Nave struct {
	PuertoDestino string `json:"puertoDestino"`
	PuertoLocal   string `json:"puertoLocal"`
	Numero        string `json:"numero"`
}

var canal = make(chan string)

//var inicio int

func recibir_mensaje(puerto_Local string) {
	confirmacion := <-canal
	fmt.Println(confirmacion)
	ln, _ := net.Listen("tcp", puerto_Local)
	defer ln.Close()
	conn, _ := ln.Accept()
	defer conn.Close()
	mensaje := bufio.NewReader(conn)
	resultado, _ := mensaje.ReadString('\n')
	if resultado == "Acceso Confirmado" {
		fmt.Println(resultado)
	} else {
		fmt.Println(resultado)
		fmt.Println("Destruido!")
		os.Exit(1)
	}
	//inicio = 1
	//inicio <- "Inicio"
}

func enviar_informacion(conn net.Conn, nave_to_json string) {
	defer conn.Close()
	fmt.Fprintf(conn, nave_to_json)
	fmt.Println(nave_to_json)
	canal <- "Terminado"
}

func main() {
	//inicio = 1
	//inicio <- "Inicio"
	// connect to this socket

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Ingrese el puerto local: ")
		puerto_Local, _ := reader.ReadString('\n')
		puerto_Local = strings.TrimSpace(puerto_Local)
		puerto_Local = ":" + puerto_Local
		//fmt.Print(puerto_Local)

		reader1 := bufio.NewReader(os.Stdin)
		fmt.Print("Ingrese el puerto destino: ")
		puerto_Destino, _ := reader1.ReadString('\n')
		puerto_Destino = strings.TrimSpace(puerto_Destino)
		puerto_Destino = ":" + puerto_Destino

		go recibir_mensaje(puerto_Local)
		//fmt.Print(puerto_Destino)

		/*
			reader_A := bufio.NewReader(os.Stdin)
			fmt.Print("Ingrese el apellido a enviar: ")
			apellido, _ := reader_A.ReadString('\n')
		*/
		//for inicio == 0 {
		//}
		//inicio = 0
		reader_B := bufio.NewReader(os.Stdin)
		fmt.Print("Ingrese el Numero a enviar: ")
		numero, _ := reader_B.ReadString('\n')
		numero = strings.TrimSpace(numero)

		objNave := Nave{
			PuertoDestino: puerto_Destino,
			PuertoLocal:   puerto_Local,
			Numero:        numero,
		}

		conn, _ := net.Dial("tcp", puerto_Destino)
		nave_to_json, _ := json.Marshal(objNave)
		enviar_informacion(conn, string(nave_to_json))
	}
}
