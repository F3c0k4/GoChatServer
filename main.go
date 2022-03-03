package main

import (
	"log"
	"net"
)

const PORT = "8080"

func main() {
	s := newServer()
	go s.execCommands()
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
		msg := "\n\nWelcome to the chat server"
		msg += "\nCommands: /nick new_nickname - Change your nickname"
		c := s.newClient(connection)

		c.sendMessage(msg)
		go c.receiveMessage()
	}

}
