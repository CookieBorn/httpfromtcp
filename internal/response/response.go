package response

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/CookieBorn/httpfromtcp/internal/headers"
)

type ResponseCode int

const (
	code200 ResponseCode = iota
	code400
	code500
)

type WriterState int

const (
	WriterStart WriterState = iota
	StatusLineDone
	HeadersDone
	BodyDone
)

type Writer struct {
	Connection   io.Writer
	State        WriterState
	ResponseCode ResponseCode
	Headers      headers.Headers
}

func (w *Writer) WriteStatusLine(statusCode ResponseCode) error {
	if w.State != WriterStart {
		return fmt.Errorf("Writer in wrong state")
	}
	switch statusCode {
	case code200:
		_, err := w.Connection.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
		w.State = StatusLineDone
		return nil
	case code400:
		_, err := w.Connection.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
		w.State = StatusLineDone
		return nil
	case code500:
		_, err := w.Connection.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
		w.State = StatusLineDone
		return nil
	default:
		return nil
	}
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != StatusLineDone {
		return fmt.Errorf("Writer in wrong state")
	}
	headString := ""
	for key, value := range headers {
		headString += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	headString += "\r\n"
	_, err := w.Connection.Write([]byte(headString))
	if err != nil {
		return err
	}
	w.State = HeadersDone
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.State != HeadersDone {
		return 0, fmt.Errorf("Writer in wrong state")
	}
	i, err := w.Connection.Write(p)
	if err != nil {
		return 0, err
	}
	w.State = BodyDone
	return i, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkHeader := fmt.Sprintf("%x\r\n", len(p))
	n1, err := w.Connection.Write([]byte(chunkHeader))
	if err != nil {
		return n1, err
	}
	n2, err := w.Connection.Write(p)
	if err != nil {
		return n1 + n2, err
	}
	n3, err := w.Connection.Write([]byte("\r\n"))
	if err != nil {
		return n1 + n2 + n3, err
	}

	return n1 + n2 + n3, nil
}
func (w *Writer) WriteChunkedBodyDone(h headers.Headers) (int, error) {
	if h == nil {
		i, err := w.Connection.Write([]byte("0\r\n\r\n"))
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	i1, err := w.Connection.Write([]byte("0\r\n"))
	if err != nil {
		return 0, err
	}
	i2, err := w.WriteTrailers(h)
	if err != nil {
		return i1, err
	}
	i3, err := w.Connection.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
	return i1 + i2 + i3, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) (int, error) {
	keys, ok := w.Headers.Get("Trailer")
	if !ok {
		return 0, fmt.Errorf("Missing Trailer Headers")
	}
	keysSplit := strings.Split(keys, ", ")
	trailerString := ""
	for _, key := range keysSplit {
		value, ok := h.Get(key)
		if !ok {
			return 0, fmt.Errorf("Missing Trailers in Header")
		}
		trailerString += key + ": " + value + "\r\n"
	}
	i, err := w.Connection.Write([]byte(trailerString))
	if err != nil {
		return i, err
	}
	return i, nil
}

func WriteStatusLine(w io.Writer, statusCode ResponseCode) error {
	switch statusCode {
	case code200:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
		return nil
	case code400:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
		return nil
	case code500:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
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
	Head["Content-Type"] = "text/html"
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
