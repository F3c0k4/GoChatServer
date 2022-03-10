package main

import (
	"log"
	"net"
)

const PORT = "8080"

func main() {
	var handler dbHandler

	err := handler.initDatabase()
	handleFatalError(err, "Unable to initiate database connection.")
	handler.pullClients()

	server := newServer(&handler)
	go server.execCommands()

	listener, err := net.Listen("tcp", ":"+PORT)
	handleFatalError(err, "Unable to start server. ")

	log.Printf("Started server on port %s", PORT)
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		handleFatalError(err, "Unable to accept connection.")
		server.newClient(connection)
	}
}

func handleFatalError(err error, desc string) {
	if err != nil {
		log.Fatalf("%s %s", desc, err)
	}
}
