package server

import (
	"net"
)

type server interface {
	listen() error
	accepter() error
	handleConnection(conn net.Conn) error
	shutdown()
}
