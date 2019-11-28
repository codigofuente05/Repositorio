package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

var remotehost string

type Message struct {
	Mensaje string `json:"mensaje"`
	Numero  string `json:"numero"`
}

func main() {
	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Remote host: ")
	remotehost, _ = gin.ReadString('\n')
	remotehost = strings.TrimSpace(remotehost)
	for {
		fmt.Print("Enter number: ")
		str, _ := gin.ReadString('\n')
		mensaje := strings.TrimSpace(str)
		message := Message{
			Mensaje: mensaje,
			Numero:  "aea",
		}
		crear_json, _ := json.Marshal(message)

		// Convertimos los datos(bytes) en una cadena e imprimimos el contenido.
		convertir_a_cadena := string(crear_json)
		fmt.Println(convertir_a_cadena)
		send(convertir_a_cadena)
	}
}

func send(num string) {
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()
	fmt.Fprintf(conn, "%s\n", num)
}
