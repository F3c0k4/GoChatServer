package main

import (
	"net"
	"strings"
	"time"
)

type server struct {
	commands        chan command
	current_clients []client
	all_clients     []db_client
	new_clients     []db_client
}

func newServer(db_clients *[]db_client) *server {
	return &server{
		commands: make(chan command),
	}
}

func (s *server) newClient(conn net.Conn) client {
	var alreadyExists bool
	client_ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	nickname := "Anonymous"
	for _, cli := range s.all_clients {
		if cli.ip == client_ip {
			nickname = cli.nickname
			alreadyExists = true
		}
	}

	c := client{
		conn:     conn,
		nickname: nickname,
		cmd:      s.commands,
	}
	s.current_clients = append(s.current_clients, c)
	if !alreadyExists {
		s.new_clients = append(s.new_clients, db_client{
			ip:       client_ip,
			nickname: nickname,
		})
	}
	return c
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
	for _, c := range s.current_clients {
		c.sendMessage("\n" + t + " " + author + ": " + msg)
	}
}
