package server

import (
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	State  bool
	Listen net.Listener
}

func Serve(port int) (*Server, error) {
	NewServer := Server{
		State: false,
	}
	listn, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		return &NewServer, err
	}
	NewServer.Listen = listn
	NewServer.State = true
	go NewServer.listen()
	return &NewServer, nil
}

func (s *Server) Close() error {
	err := s.Listen.Close()
	if err != nil {
		return err
	}
	s.State = false
	return nil
}

func (s *Server) listen() {
	if s.State {
		con, err := s.Listen.Accept()
		if err != nil {
			fmt.Printf("Listen error: %s\n", err)
		}
		s.handle(con)
	}
}

func (s *Server) handle(conn net.Conn) {
	_, err := conn.Write([]byte("HTTP/1.1 200 OK\nContent-Type: text/plain\n\nHello World!"))
	if err != nil {
		fmt.Printf("Write error: %s\n", err)
	}
	err = conn.Close()
	if err != nil {
		fmt.Printf("handle error: %s\n", err)
	}
}
