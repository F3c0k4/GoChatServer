package main

type CommandId int

const (
	CMD_NICK CommandId = iota
	CMD_EXIT
	CMD_BROADCAST
)

type command struct {
	cmd_id CommandId
	args   []string
	client *client
}
