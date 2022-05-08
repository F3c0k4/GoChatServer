package main

import (
	"GoChatServer/pkg/server"
	"log"
	"net"
)

const PORT = "8080"

func main() {
	var handler server.DbHandler

	// starting database and loading clients
	err := handler.InitDatabase()
	handleFatalError(err, "Unable to initiate database connection.")
	handler.PullClients()

	// creating server and starting goroutine to exec commands
	server := server.NewServer(&handler)
	go server.ExecCommands()

	//starting server
	listener, err := net.Listen("tcp", ":"+PORT)
	handleFatalError(err, "Unable to start server. ")
	defer listener.Close()
	log.Printf("Started server on port %s", PORT)

	// handling new connections
	for {
		connection, err := listener.Accept()
		handleFatalError(err, "Unable to accept connection.")
		server.NewClient(connection)
	}
}

func handleFatalError(err error, desc string) {
	if err != nil {
		log.Fatalf("%s %v", desc, err)
	}
}
