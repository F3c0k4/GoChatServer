package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nickname string
	cmd      chan<- command
}

func (c *client) sendMessage(msg string) {
	c.conn.Write([]byte(msg))
}

func (c *client) receiveMessage() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')

		if err != nil {
			c.sendMessage(fmt.Sprintf("Error while trying to receive your message: %s", err.Error()))
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")

		if args[0] == "/nick" {
			c.cmd <- command{
				cmd_id: CMD_NICK,
				args:   args[1:],
				client: c,
			}
		} else {
			c.cmd <- command{
				cmd_id: CMD_BROADCAST,
				args:   args,
				client: c,
			}
		}
	}
}
