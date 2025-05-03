package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	strRead := string(data)
	i := strings.Index(strRead, "\r\n")
	if i < 0 {
		return 0, false, nil
	} else if i == 0 {
		return 2, true, nil
	}
	headTrim := strings.TrimSpace(strRead[:i])
	headParts := strings.Split(headTrim, ":")
	headParts[1] = strings.TrimSpace(headParts[1])
	for _, headPart := range headParts {
		if headPart != strings.TrimSpace(headPart) {
			return 0, false, fmt.Errorf("Header not formated correctly")
		}
	}
	if len(headParts) > 2 {
		h[headParts[0]] = headParts[1] + ":" + headParts[2]
	} else {
		h[headParts[0]] = headParts[1]
	}

	return len(headTrim) + 2, false, nil
}
