package main

import (
	"GoChatServer/pkg/server"
	"log"
	"net"
)

const PORT = "8080"

func main() {
	var handler server.DbHandler

	err := handler.InitDatabase()
	handleFatalError(err, "Unable to initiate database connection.")
	handler.PullClients()

	server := server.NewServer(&handler)
	go server.ExecCommands()

	listener, err := net.Listen("tcp", ":"+PORT)
	handleFatalError(err, "Unable to start server. ")

	log.Printf("Started server on port %s", PORT)
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		handleFatalError(err, "Unable to accept connection.")
		server.NewClient(connection)
	}
}

func handleFatalError(err error, desc string) {
	if err != nil {
		log.Fatalf("%s %s", desc, err)
	}
}
