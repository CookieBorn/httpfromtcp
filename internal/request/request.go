package request

import (
	"fmt"
	"io"
	"strings"

	"github.com/CookieBorn/httpfromtcp/internal/headers"
)

// State. 0: Parse Requests, 1: Parse Headers, 2: done

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
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

	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}

		readBuffer = append(readBuffer, buf[:n]...)
		bytesProcessed, parseErr := req.parse(readBuffer)
		if parseErr != nil {
			return nil, parseErr
		}

		// Trim the buffer to remove processed bytes
		readBuffer = readBuffer[bytesProcessed:]

		// If we've reached the end of headers or EOF, we're done
		if req.State == 2 || (err == io.EOF && n == 0) {
			break
		}
	}

	return req, nil
}

func parseRequestLine(res []byte) (int, RequestLine, error) {
	returnStruct := RequestLine{}
	strRead := string(res)
	i := strings.Index(strRead, "\r\n")
	if i < 0 {
		return 0, returnStruct, nil
	}
	parts := strings.Split(strRead, "\r\n")
	reqLine := parts[0]
	reqLineParts := strings.Split(reqLine, " ")
	if len(reqLineParts) != 3 {
		return i, returnStruct, fmt.Errorf("Not enough parts in request line")
	}
	if reqLineParts[0] != strings.ToUpper(reqLineParts[0]) {
		return i, returnStruct, fmt.Errorf("Method incorect in request line")
	}
	if reqLineParts[2] != "HTTP/1.1" {
		return i, returnStruct, fmt.Errorf("only HTTP/1.1 accepted")
	}
	returnStruct.HttpVersion = "1.1"
	returnStruct.Method = reqLineParts[0]
	returnStruct.RequestTarget = reqLineParts[1]
	return i, returnStruct, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != 2 {
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
		i, reqLine, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if i == 0 {
			// Need more data
			return 0, nil
		}
		r.RequestLine = reqLine
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
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unexpected state: %v", r.State)
	}
}
