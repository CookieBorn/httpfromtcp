package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Errorf("Open error: %v\n", err)
		return
	}
	chanell := getLinesChannel(file)
	for Line := range chanell {
		if Line == "ay?ay?" {
			break
		} else {
			fmt.Printf("read: %s\n", Line)
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
