package main

import (
	"log"
	"net"
)

const PORT = "8080"

func main() {
	listener, err := net.Listen("tcp", ":"+PORT)

	if err != nil {
		log.Fatalf("Unable to start server. Error: %s", err.Error())
	}

	log.Printf("Started server on port %s", PORT)
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("Unable to accept connection. Error: %s", err.Error())
		}

		log.Printf("Connections IP address is : %s", connection.RemoteAddr().(*net.TCPAddr).IP.String())
	}

}
