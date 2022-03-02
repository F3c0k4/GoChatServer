package main

import (
	"fmt"
	"net"
)

type server struct {
	commands chan command
}

func newServer() *server {
	return &server{
		commands: make(chan command),
	}
}

func (s *server) newClient(conn net.Conn) *client {
	c := client{
		conn:     conn,
		nickname: "Anonymous",
		cmd:      s.commands,
	}

	return &c
}

func (s *server) execCommands() {
	for c := range s.commands {
		switch c.cmd_id {
		case CMD_NICK:
			s.changeNick(c.client, c.args)
		case CMD_BROADCAST:
			s.broadcastMessage(c.args)
		}
	}
}

func (s *server) changeNick(c *client, args []string) {
	c.nickname = args[0]
	c.sendMessage(fmt.Sprintf("\nYour new nickname is %s", c.nickname))
}

func (s *server) broadcastMessage(args []string) {
	//TODO
}
