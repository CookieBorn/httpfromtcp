package main

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/CookieBorn/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		fmt.Errorf("Listener error: %v\n", err)
		return
	}
	defer listener.Close()
	defer fmt.Print("Channel has been closed\n")
	for true {
		con, err := listener.Accept()
		if err != nil {
			fmt.Errorf("Connection error: %v\n", err)
			return
		}
		fmt.Print("A connection has been established\n")
		req, err := request.RequestFromReader(con)
		if err != nil {
			fmt.Errorf("Request error: %v\n", err)
			return
		}
		fmt.Print("Request line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for head := range req.Headers {
			fmt.Printf("- %s: %s\n", head, req.Headers[head])
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	chanString := make(chan string)
	b := make([]byte, 8)
	Line := ""
	go func() {
		for 1 > 0 {
			n, err := f.Read(b)
			if n > 0 {
				Line += string(b[:n])
				for {
					i := strings.Index(Line, "\n")
					if i < 0 {
						break
					}
					sendLine := Line[:i]
					chanString <- sendLine
					Line = Line[i+1:]
				}
			}
			if err != nil {
				if Line != "" {
					chanString <- Line
				}
				f.Close()
				close(chanString)
				return
			}
		}
	}()
	return chanString
}
