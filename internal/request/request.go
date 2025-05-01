package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	reqLine, err := parseRequestLine(reBytes)
	if err != nil {
		return nil, err
	}
	req := Request{
		RequestLine: reqLine,
	}
	return &req, nil
}

func parseRequestLine(res []byte) (RequestLine, error) {
	returnStruct := RequestLine{}
	strRead := string(res)
	parts := strings.Split(strRead, "\r\n")
	reqLine := parts[0]
	reqLineParts := strings.Split(reqLine, " ")
	if len(reqLineParts) != 3 {
		return returnStruct, fmt.Errorf("Not enough parts in request line")
	}
	if reqLineParts[0] != strings.ToUpper(reqLineParts[0]) {
		return returnStruct, fmt.Errorf("Method incorect in request line")
	}
	if reqLineParts[2] != "HTTP/1.1" {
		return returnStruct, fmt.Errorf("only HTTP/1.1 accepted")
	}
	returnStruct.HttpVersion = "1.1"
	returnStruct.Method = reqLineParts[0]
	returnStruct.RequestTarget = reqLineParts[1]
	return returnStruct, nil
}
