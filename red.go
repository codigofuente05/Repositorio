package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var lista []int

var destroyers []string
var servidores_naves []string

var puerto_Local_Global string
var puerto_Dial_Global string

var canal = make(chan string)

type Nave struct {
	PuertoDestino string `json:"puertoDestino"`
	PuertoLocal   string `json:"puertoLocal"`
	Numero        string `json:"numero"`
}
type Destroyer struct {
	PuertoDial   string `json:"puertoDial"`
	PuertoListen string `json:"puertoListen"`
}

func cliente_para_nave(mensaje, nave_destino string) {
	fmt.Println("Notificar mensaje a nave espacial...")
	ln, _ := net.Dial("tcp", nave_destino)
	defer ln.Close()
	fmt.Fprintf(ln, mensaje)
}
func generarLista() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := r1.Intn(10)
	for i := 0; i < n; i++ {
		lista = append(lista, r1.Intn(100))
	}
}
func verificar_numero(numero int) bool {
	for i := 0; i < len(lista); i++ {
		if lista[i] == numero {
			return true
		}
	}
	return false
}
func enviar_informacion(conn net.Conn, nave_to_json string) {
	defer conn.Close()
	fmt.Fprintf(conn, nave_to_json)
	fmt.Println(nave_to_json)
	//canal <- "Terminado"
}
func verificar_numero_en_otros_destroyer(objNave Nave) {
	for i, _ := range destroyers {
		conn, _ := net.Dial("tcp", destroyers[i])
		nave_to_json, _ := json.Marshal(objNave)
		enviar_informacion(conn, string(nave_to_json))
	}
}
func leer_informacion(conn net.Conn) {
	defer conn.Close()
	result := bufio.NewReader(conn)
	message, _ := result.ReadString('\n')
	var objNave Nave
	json.Unmarshal([]byte(message), &objNave)
	numero_to_int, _ := strconv.Atoi(objNave.Numero)
	fmt.Printf("Numero: %d\n", numero_to_int)

	var mensaje string
	if verificar_numero(numero_to_int) {
		mensaje = "Acceso Confirmado"
	} else {
		mensaje = "Acceso Denegado"
		//verificar_numero_en_otros_destroyer(objNave)
	}
	fmt.Println(mensaje)
	cliente_para_nave(mensaje, objNave.PuertoLocal)
}
func servidor_para_nave(puerto_Local string) {
	ln, _ := net.Listen("tcp", puerto_Local)
	defer ln.Close()
	for {
		fmt.Println("Esperando...")
		conn, err1 := ln.Accept()
		if err1 != nil {
			fmt.Println("Error ", err1)
		} else {
			fmt.Println("Llego!")
			//canal <- "Inicio"
			go leer_informacion(conn)
			//<-canal
		}
	}
}
func conenctarse_a_destroyer_2(objDestroyer Destroyer, destino string) {
	ln, _ := net.Dial("tcp", destino)
	defer ln.Close()
	destroyer_to_json, _ := json.Marshal(objDestroyer)
	fmt.Fprintf(ln, string(destroyer_to_json))
}
func enviarAOtros(objDestroyer Destroyer) {
	for i, _ := range destroyers {
		for j, _ := range destroyers {
			if destroyers[j] != destroyers[i] {
				objDestroyer2 := Destroyer{
					PuertoDial:   "",
					PuertoListen: destroyers[j],
				}
				conenctarse_a_destroyer_2(objDestroyer2, destroyers[i])
			}
		}
	}
}
func verificar(PuertoListen string) bool {
	for index := 0; index < len(destroyers); index++ {
		if destroyers[index] == PuertoListen {
			return true
		}
	}
	return false
}
func almacenar_informacion_otros_destroyer(conn net.Conn) {
	defer conn.Close()
	result := bufio.NewReader(conn)
	message, _ := result.ReadString('\n')
	var objDestroyer Destroyer
	json.Unmarshal([]byte(message), &objDestroyer)
	if verificar(objDestroyer.PuertoListen) == false {
		//servidores_naves = append(servidores_naves, objDestroyer.PuertoDial)
		destroyers = append(destroyers, objDestroyer.PuertoListen)
		enviarAOtros(objDestroyer)
	}
	//fmt.Println(servidores_naves)
	canal <- "Inicio"
}

func servidor_para_destroyer(puerto_Local1 string) {
	ln, _ := net.Listen("tcp", puerto_Local1)
	defer ln.Close()
	for {
		//fmt.Println("Esperando...")
		conn, err1 := ln.Accept()
		if err1 != nil {
			fmt.Println("Error ", err1)
		} else {
			//fmt.Println("Llego!")
			//canal <- "Inicio"
			go almacenar_informacion_otros_destroyer(conn)
			//<-canal
			go func() {
				<-canal
				fmt.Println(destroyers)
			}()
		}
	}
}
func conenctarse_a_destroyer(objDestroyer Destroyer) {
	ln, _ := net.Dial("tcp", puerto_Dial_Global)
	defer ln.Close()
	destroyer_to_json, _ := json.Marshal(objDestroyer)
	fmt.Fprintf(ln, string(destroyer_to_json))
}
func main() {
	generarLista()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto local (Servidor Para Nave): ")
	puerto_Local, _ := reader.ReadString('\n')
	puerto_Local = strings.TrimSpace(puerto_Local)
	puerto_Local = ":" + puerto_Local

	//servidores_naves = append(servidores_naves, puerto_Local)

	reader1 := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto local (Servidor Para Destroyer): ")
	puerto_Local1, _ := reader1.ReadString('\n')
	puerto_Local1 = strings.TrimSpace(puerto_Local1)
	puerto_Local1 = ":" + puerto_Local1
	puerto_Local_Global = puerto_Local1

	reader2 := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto remoto (Conectarse a otro Destroyer): ")
	puerto2, _ := reader2.ReadString('\n')
	puerto2 = strings.TrimSpace(puerto2)

	go servidor_para_nave(puerto_Local)
	go servidor_para_destroyer(puerto_Local1)

	if len(puerto2) > 0 {
		puerto2 = ":" + puerto2
		puerto_Dial_Global = puerto2
		objDestroyer := Destroyer{
			PuertoDial:   puerto_Local,
			PuertoListen: puerto_Local1,
		}
		destroyers = append(destroyers, puerto2)
		go conenctarse_a_destroyer(objDestroyer)
	}

	fmt.Print(lista)
	for {
	}
}

//8081 - Recibe info de nave
//8082 - Envia info a nave

//9081 - Recibe info de Destroyer
//9082 - Envia info a Destroyer
