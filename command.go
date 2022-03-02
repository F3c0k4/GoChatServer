package main

type CommandId int

const (
	CMD_NICK CommandId = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_EXIT
	CMD_BROADCAST
)

type command struct {
	cmd_id CommandId
	args   []string
	client *client
}
