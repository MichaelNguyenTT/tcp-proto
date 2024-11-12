package main

import (
	"tcpserver/internal/server"
)

func main() {
	server := server.NewServer(":8080")
	server.Run()
}
