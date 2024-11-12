package server

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type tcpServer struct {
	addr     string
	listener net.Listener

	ready chan bool
	done  chan bool
}

func NewServer(address string) *tcpServer {
	return &tcpServer{
		addr:  address,
		ready: make(chan bool),
		done:  make(chan bool),
	}
}

func (s *tcpServer) listen() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.listener = ln

	close(s.ready)

	return s.accepter()
}

// accept the loop
func (s *tcpServer) accepter() error {
	defer s.listener.Close()

	for {
		select {
		case <-s.done:
			slog.Info("graceful shutdown")
			return nil

		default:
			conn, err := s.listener.Accept()
			if err != nil {
				return fmt.Errorf("failed to accept connection: %v", err)
			}

			slog.Info("accepted connection", "addr", conn.RemoteAddr())
			go s.handleConnection(conn)
		}
	}
}

// read data from client connection
func (s *tcpServer) handleConnection(conn net.Conn) error {
	defer conn.Close()

	//TODO: create your own data processing protocol
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				slog.Info("connection closed", "addr", conn.RemoteAddr())
			} else {
				slog.Error("failed to read from connection", "err", err)
			}
			return err
		}

		msg := buf[:n]

		slog.Info("received from connection", "from", conn.RemoteAddr().String(), "msg", msg)
	}
}

func (s *tcpServer) shutdown() {
	close(s.done)
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *tcpServer) Run() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := s.listen()
		if err != nil {
			slog.Error("failed to start server", "msg", err)
			os.Exit(1)
			return
		}
	}()

	<-s.ready
	slog.Info("Server ready to accept connections:", "port", s.addr)

	<-sig
	s.shutdown()
	slog.Info("Server shutdown complete")
}
