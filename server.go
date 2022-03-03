package main

import (
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
			s.broadcastMessage(c.client.nickname, strings.Join(c.args, " "))
		}
	}
}

func (s *server) changeNick(c *client, args []string) {
	if len(args) > 0 {
		msg := c.nickname + " changed their nickname to " + args[0]
		s.broadcastMessage("Server", msg)
		c.nickname = args[0]
	} else {
		msg := "\nName argument missing, please try again. Usage /nick new_nickname "
		c.sendMessage(msg)
	}

}

func (s *server) broadcastMessage(author string, msg string) {
	t := time.Now().Format("15:04:05")
	for _, c := range s.clients {
		c.sendMessage("\n" + t + " " + author + ": " + msg)
	}
}
