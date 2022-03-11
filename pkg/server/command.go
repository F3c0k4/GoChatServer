package server

type CommandId int

// Command ids used for recognition of server commands
const (
	CMD_NICK CommandId = iota
	CMD_EXIT
	CMD_BROADCAST
	CMD_REMOVE_CLIENT
)

// command stores information about the service to be provided
// by the server to client
type command struct {
	cmdId  CommandId
	args   []string
	client *client
}
