package headers

import (
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	strRead := string(data)
	i := strings.Index(strRead, "\r\n")
	if i < 0 {
		return 0, false, nil // Need more data
	} else if i == 0 {
		return 2, true, nil
	}
	headerLine := strRead[:i]
	headParts := strings.SplitN(headerLine, ":", 2)
	if len(headParts) != 2 {
		return 0, false, fmt.Errorf("malformed header: %s", headerLine)
	}

	name := strings.ToLower(headParts[0])
	if name != strings.TrimRight(name, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", name)
	}
	value := strings.TrimSpace(headParts[1])

	if !ValidnameCheck(strings.ToLower(headParts[0])) {
		return 0, false, fmt.Errorf("Incorrect characters in runes")
	}
	if _, ok := h[name]; ok {
		h[name] += ", " + value
	} else {
		h[name] = value
	}
	return len(headerLine) + 2, false, nil
}

func ValidnameCheck(name string) bool {
	allowed := []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}
	NotSpecial := []rune{}
	for _, letter := range name {
		found := false
		for _, all := range allowed {
			if all == letter {
				found = true
			}
		}
		if found == false {
			NotSpecial = append(NotSpecial, letter)
		}
	}
	for _, letter := range NotSpecial {
		if !unicode.IsLetter(letter) && !unicode.IsNumber(letter) {
			return false
		}
	}

	return true
}
