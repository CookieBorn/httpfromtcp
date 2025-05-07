package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/CookieBorn/httpfromtcp/internal/headers"
)

const crlf = "\r\n"

// State. 0: Parse Requests, 1: Parse Headers, 2: Parse Body, 3: Done

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	State       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := &Request{
		State:   0, // requestStateParsingRequestLine
		Headers: headers.NewHeaders(),
	}
	buf := make([]byte, 4096)
	readBuffer := []byte{}
	for req.State != 3 {
		n, err := reader.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != 3 {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.State, n)
				}
				break
			}
			return nil, err
		}
		readBuffer = append(readBuffer, buf[:n]...)
		bytesProcessed, parseErr := req.parse(readBuffer)
		if parseErr != nil {
			return nil, parseErr

		}
		readBuffer = readBuffer[bytesProcessed:]
	}

	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != 3 {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case 0:
		// Your request line parsing logic
		reqLine, i, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if i == 0 {
			// Need more data
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.State = 1
		return i, nil
	case 1:
		i, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = 2
		}
		return i, nil
	case 2:
		length, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.State = 3
			return 0, nil
		}
		lengthInt, err := strconv.Atoi(length)
		if err != nil {
			return 0, fmt.Errorf("Content-Length jeader malformed")
		}
		r.Body = data
		if len(r.Body) > lengthInt {
			return len(r.Body), fmt.Errorf("Data length: %d, Content-length: %d. Should be equal\n", len(data), lengthInt)
		} else if len(r.Body) == lengthInt {
			r.State = 3
			return lengthInt, nil
		}
		return 0, nil
	case 3:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unexpected state: %v", r.State)
	}
}
