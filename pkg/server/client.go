package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

// client stores the connection handler of a client,
// their nickname and the commands they send to the server
type client struct {
	conn     net.Conn
	nickname string
	cmd      chan<- command
}

// sendMessage sends the msg parameter to a client
func (c *client) sendMessage(msg string) {
	c.conn.Write([]byte(msg))
}

// receiveMessage listens for a message from a client in an endless loop
// and determines the command to be executed based on the content of the message
func (c *client) receiveMessage() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')

		if err == io.EOF {
			c.conn.Close()
			c.cmd <- command{
				cmdId:  CMD_REMOVE_CLIENT,
				args:   nil,
				client: c,
			}
			return
		} else if err != nil {
			c.sendMessage(fmt.Sprintf("Error while trying to receive your message: %s", err.Error()))
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")

		if args[0] == "/nick" {
			c.cmd <- command{
				cmdId:  CMD_NICK,
				args:   args[1:],
				client: c,
			}
		} else {
			c.cmd <- command{
				cmdId:  CMD_BROADCAST,
				args:   args,
				client: c,
			}
		}
	}
}
