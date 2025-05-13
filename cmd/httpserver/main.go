package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CookieBorn/httpfromtcp/internal/request"
	"github.com/CookieBorn/httpfromtcp/internal/response"
	"github.com/CookieBorn/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, SecondHandle)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func FirstHandle(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			Code:     1,
			ErrorMSG: "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			Code:     2,
			ErrorMSG: "Woopsie, my bad\n",
		}
	default:
		_, err := w.Write([]byte("All good, frfr\n"))
		if err != nil {
			return &server.HandlerError{
				Code:     2,
				ErrorMSG: "Failed to write response",
			}
		}
		return nil
	}
}

func SecondHandle(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		_, err := w.Connection.Write([]byte("<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>"))
		w.ResponseCode = 1
		if err != nil {
			fmt.Printf("Write handle error: %s", err)
			return
		}
		return
	case "/myproblem":
		_, err := w.Connection.Write([]byte("<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>"))
		w.ResponseCode = 2
		if err != nil {
			fmt.Printf("Write handle error: %s", err)
			return
		}
		return
	default:
		_, err := w.Connection.Write([]byte("<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>\n"))
		w.ResponseCode = 0
		if err != nil {
			fmt.Printf("Write handle error: %s", err)
			return
		}
		return
	}
}
