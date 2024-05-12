package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

type Client chan<- string

var (
	incomingClients = make(chan Client)
	leavingClients  = make(chan Client)
	messages        = make(chan string)
)

var (
	host = flag.String("h", "localhost", "host")
	port = flag.Int("p", 3090, "port")
)

// Client1 -> Server -> HandleConnection(Client1)

func HandleUserConnection(conn net.Conn) {
	defer conn.Close()

	// Canal del cliente (específicamente de el)
	clientMessages := make(chan string)
	go MessageWrite(conn, clientMessages)

	// Generamos un nombre para el cliente
	clientName := conn.RemoteAddr().String()

	// Mensaje para el nuevo cliente conectado al servidor.
	clientMessages <- fmt.Sprintf("Welcome to the server, %s\n!", clientName)

	// Mensaje de bienvenida en el chat global
	messages <- fmt.Sprintln("New client connected! Say hi to", clientName)

	// Agregar el cliente a la lista de clientes conectados
	incomingClients <- clientMessages

	// Lee los mensajes del ciente
	// Si la conexión se cierra, el bucle se rompe
	inputMessage := bufio.NewScanner(conn)
	for inputMessage.Scan() {
		messages <- fmt.Sprintf("%s: %s", clientName, inputMessage.Text())
	}

	// Agrega el cliente a la lista de clientes saliendo del chat
	leavingClients <- clientMessages

	// Mensaje global anunciando que el usuario se desconectó
	messages <- fmt.Sprintln(clientName, "has leaved the chat.")
}

// Recibe mensajes del canal y los escribe en el cliente
func MessageWrite(conn net.Conn, clientMessages <-chan string) {
	for message := range clientMessages {
		fmt.Fprintln(conn, message)
	}
}

func Broadcast() {
	clients := make(map[Client]bool)
	for {
		select {
		case message := <-messages:
			for client := range clients {
				client <- message
			}

		case newClient := <-incomingClients:
			clients[newClient] = true

		case leavingClient := <-leavingClients:
			delete(clients, leavingClient)
			close(leavingClient)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))

	if err != nil {
		log.Fatal(err)
	}

	go Broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go HandleUserConnection(conn)
	}

}
