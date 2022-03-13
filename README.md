# GoChatServer is a TCP chat server written in Go.
Prerequisite:
 - Replace the data in the assets/credentials.env file with the credentials to your Postgresql database

Usage:
 - Run cmd/main.go
 - You can connect to the server using telnet through port 8080
 - You can set your nickname with "\nick nickname" command
 - Other messages are considered chat messages and will be broadcast to all online users
 - If you change your nickname, the server will remember it when you reconnect, even if the server restarts in the meantime
