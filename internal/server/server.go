package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/CookieBorn/httpfromtcp/internal/response"
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

	var resCode response.ResponseCode = 0
	err := response.WriteStatusLine(conn, resCode)
	if err != nil {
		fmt.Printf("Write status line error: %s\n", err)
	}

	head := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, head)
	if err != nil {
		fmt.Printf("Write header error: %s\n", err)
	}

	err = conn.Close()
	if err != nil {
		fmt.Printf("handle error: %s\n", err)
	}
}
