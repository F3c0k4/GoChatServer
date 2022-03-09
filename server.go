package main

import (
	"net"
	"strings"
	"time"
)

// server stores the commands to be executed,
// the clients currently connected to the server
// and the handler of the database
type server struct {
	commands       chan command
	currentClients []client
	handler        *dbHandler
}

// newServer creates and returns a pointer to a new server
func newServer(handler *dbHandler) *server {

	return &server{
		commands: make(chan command),
		handler:  handler,
	}
}

// newClient adds a client to the currentClients list and
// updates the database if the client is new
func (s *server) newClient(conn net.Conn) {
	var alreadyExists bool
	var msg string
	client_ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	nickname := "Anonymous"

	for _, cli := range s.handler.dbClients {
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
	s.currentClients = append(s.currentClients, c)
	if !alreadyExists {
		new_client := dbClient{
			ip:       client_ip,
			nickname: nickname,
		}
		s.handler.addClient(new_client)
		msg = "\n\nWelcome to the chat server"
		msg += "\nCommands: /nick new_nickname - Change your nickname"
	} else {
		msg = "\n\nWelcome back to the chat server, " + nickname
		msg += "\nCommands: /nick new_nickname - Change your nickname"
	}

	c.sendMessage(msg)
	go c.receiveMessage()
}

// execCommands executes the commands previously saved by the clients
func (s *server) execCommands() {
	for c := range s.commands {
		switch c.cmdId {
		case CMD_NICK:
			s.changeNick(c.client, c.args)
		case CMD_BROADCAST:
			s.broadcastMessage(c.client.nickname, strings.Join(c.args, " "))
		}
	}
}

// changeNick changes the nickname of a client and also updates the
// database
func (s *server) changeNick(c *client, args []string) {
	if len(args) > 0 {
		msg := c.nickname + " changed their nickname to " + args[0]
		s.broadcastMessage("Server", msg)
		c.nickname = args[0]
		db_cli := dbClient{
			ip:       c.conn.RemoteAddr().(*net.TCPAddr).IP.String(),
			nickname: c.nickname,
		}
		s.handler.updateClientRecord(db_cli)
	} else {
		msg := "\nName argument missing, please try again. Usage /nick new_nickname "
		c.sendMessage(msg)
	}

}

// boradcastMessage sends a message to all the currently connected
// clients
func (s *server) broadcastMessage(author string, msg string) {
	t := time.Now().Format("15:04:05")
	for _, c := range s.currentClients {
		c.sendMessage("\n" + t + " " + author + ": " + msg)
	}
}
