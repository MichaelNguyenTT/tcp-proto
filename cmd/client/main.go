package main

import (
	"fmt"
	"log/slog"
	"net"
)

const (
	DefaultHost = "localhost"
	DefaultPort = "8080"
)

type client struct {
	host string
	port string
}

func NewClient() *client {
	return &client{
		host: DefaultHost,
		port: DefaultPort,
	}
}

func (c *client) Addr() string {
	return c.host + ":" + c.port
}

func (c *client) Connect() error {
	conn, err := net.Dial("tcp", c.Addr())
	if err != nil {
		return fmt.Errorf("couldn't connect to server: %v", err)
	}

	defer conn.Close()
	slog.Info("connection sent to server", "addr", conn.RemoteAddr().String())

	sendmsg := "hello, from TCP Client"
	_, err = conn.Write([]byte(sendmsg))
	if err != nil {
		slog.Error("failed to send msg to server", "err", err)
		return err
	}

	return nil
}

func main() {
	client := NewClient()
	client.Connect()
}
