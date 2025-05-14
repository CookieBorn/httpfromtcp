package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/CookieBorn/httpfromtcp/internal/headers"
	"github.com/CookieBorn/httpfromtcp/internal/request"
	"github.com/CookieBorn/httpfromtcp/internal/response"
	"github.com/CookieBorn/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, ThirdHandle)
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

func ThirdHandle(w *response.Writer, req *request.Request) {
	if ok := strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/"); !ok {
		w.Headers = nil
		w.ResponseCode = 1
		w.Connection.Write([]byte("Missing prefix"))
		return
	}
	length := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

	res, err := http.Get("https://httpbin.org/" + length)
	if err != nil {
		w.Headers = nil
		w.ResponseCode = 2
		w.Connection.Write([]byte("Get error"))
		return
	}
	defer res.Body.Close()

	w.ResponseCode = 0

	head := headers.NewHeaders()
	for k, vList := range res.Header {
		if k != "Content-Length" {
			for _, v := range vList {
				head.Set(k, v)
			}
		}
	}
	head.Set("Transfer-Encoding", "chunked")
	head.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	w.Headers = head

	buf := make([]byte, 1024)
	resBody := ""
	for {
		n, err := res.Body.Read(buf)
		if n > 0 {
			w.WriteChunkedBody(buf[:n])
			resBody += string(buf[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			w.Headers = nil
			w.ResponseCode = 1
			w.Connection.Write([]byte("Read error"))
			return
		}
	}
	trailers := headers.NewHeaders()
	trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", sha256.Sum256([]byte(resBody))))
	trailers.Set("X-Content-Length", strconv.Itoa(len([]byte(resBody))))
	w.WriteChunkedBodyDone(trailers)

	return
}
