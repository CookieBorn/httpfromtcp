package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	State       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req := Request{
		State: 1,
	}
	var buf []byte
	for req.State == 1 {
		reBytes, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		buf = append(buf, reBytes...)
		_, err = req.parse(buf)
		if err != nil {
			return nil, err
		}
	}
	return &req, nil
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
	i, reqLine, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if i == 0 {
		r.State = 1
		return 0, nil
	}
	r.RequestLine = reqLine
	r.State = 0
	return 0, nil
}
