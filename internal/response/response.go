package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/CookieBorn/httpfromtcp/internal/headers"
)

type ResponseCode int

const (
	code200 ResponseCode = iota
	code400
	code500
)

func WriteStatusLine(w io.Writer, statusCode ResponseCode) error {
	switch statusCode {
	case code200:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
		return nil
	case code400:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\n"))
		if err != nil {
			return err
		}
		return nil
	case code500:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\n"))
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	Head := headers.NewHeaders()
	Head["Content-Length"] = strconv.Itoa(contentLen)
	Head["Content-Type"] = "text/plain"
	Head["Connection"] = "close"
	return Head
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headString := ""
	for key, value := range headers {
		headString += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	headString += "\r\n"
	_, err := w.Write([]byte(headString))
	if err != nil {
		return err
	}
	return nil
}
