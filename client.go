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
}

func (c *client) sendMessage(msg string) {
	c.conn.Write([]byte(msg))
}

func (c *client) receiveMessage() {
	msg, err := bufio.NewReader(c.conn).ReadString('\n')

	if err != nil {
		c.sendMessage(fmt.Sprintf("Error while trying to receive your message: %s", err.Error()))
	}

	msg = strings.Trim(msg, "\r\n")
	args := strings.Split(msg, " ")

	if args[0] == "/nick" {
		c.nickname = args[1]
		c.sendMessage(fmt.Sprintf("\nYour new nickname is %s", c.nickname))
	}

}
