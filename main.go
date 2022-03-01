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
}
