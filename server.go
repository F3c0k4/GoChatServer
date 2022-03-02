package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type server struct {
	commands chan command
	clients  []*client
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
	s.clients = append(s.clients, &c)
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
	t := time.Now().Format("2006-01-02 15:04:05")
	for _, c := range s.clients {
		c.sendMessage("\n" + t + " | " + c.nickname + ": " + strings.Join(args, " "))
	}
}
