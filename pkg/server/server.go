package server

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
	handler        *DbHandler
}

// NewServer creates and returns a pointer to a new server
func NewServer(handler *DbHandler) *server {

	return &server{
		commands: make(chan command),
		handler:  handler,
	}
}

// NewClient adds a client to the currentClients list and
// updates the database if the client is new
func (s *server) NewClient(conn net.Conn) {
	var alreadyExists bool
	var msg string
	client_ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	nickname := "Anonymous"

	for _, cli := range s.handler.dbClients {
		if cli.ip == client_ip {
			nickname = cli.nickname
			alreadyExists = true
			break
		}
	}
	c := client{
		conn:     conn,
		nickname: nickname,
		cmd:      s.commands,
	}
	s.currentClients = append(s.currentClients, c)
	if !alreadyExists {
		new_client := DbClient{
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

// ExecCommands executes the commands previously saved by the clients
func (s *server) ExecCommands() {

	for c := range s.commands {
		switch c.cmdId {
		case CMD_NICK:
			s.changeNick(c.client, c.args)
		case CMD_BROADCAST:
			s.broadcastMessage(c.client.nickname, strings.Join(c.args, " "))
		case CMD_REMOVE_CLIENT:
			s.removeCurrentClient(c.client)
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
		db_cli := DbClient{
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

// removeCurrentClient removes a client from the currentClients list
func (s *server) removeCurrentClient(c *client) {
	for i, cli := range s.currentClients {
		if c.conn == cli.conn {
			s.currentClients = append(s.currentClients[:i], s.currentClients[i+1:]...)
			break
		}
	}
}