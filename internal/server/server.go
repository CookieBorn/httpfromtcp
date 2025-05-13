package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/CookieBorn/httpfromtcp/internal/request"
	"github.com/CookieBorn/httpfromtcp/internal/response"
)

type Server struct {
	State       bool
	HandlerFunc Handler
	Listen      net.Listener
}

type Handler func(w *response.Writer, req *request.Request)

type HandlerError struct {
	Code     response.ResponseCode
	ErrorMSG string
}

func Serve(port int, handler Handler) (*Server, error) {
	NewServer := Server{
		State:       false,
		HandlerFunc: handler,
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
	for {
		con, err := s.Listen.Accept()
		if err != nil {
			fmt.Printf("Listen error: %s\n", err)
		}
		go s.handle(con)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("Parse Req error: %s", err)
		return
	}

	buffer := bytes.NewBuffer(nil)

	writer := response.Writer{
		Connection:   buffer,
		State:        response.WriterStart,
		ResponseCode: 0,
	}

	s.HandlerFunc(&writer, req)

	head := response.GetDefaultHeaders(buffer.Len())

	err = response.WriteStatusLine(conn, writer.ResponseCode)
	if err != nil {
		fmt.Printf("Write status line error: %s\n", err)
		return
	}

	err = response.WriteHeaders(conn, head)
	if err != nil {
		fmt.Printf("Write header error: %s\n", err)
		return
	}

	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		fmt.Printf("Write body error: %s\n", err)
		return
	}
}

func Error(w io.Writer, handError HandlerError) error {
	err := response.WriteStatusLine(w, handError.Code)
	if err != nil {
		return err
	}

	headers := response.GetDefaultHeaders(len(handError.ErrorMSG))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(handError.ErrorMSG))
	if err != nil {
		return err
	}
	return nil
}
