package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Errorf("Open error: %v\n", err)
		return
	}
	b := make([]byte, 8)
	Line := []string{}
	for 1 > 0 {
		_, err := file.Read(b)
		if err != nil {

			fmt.Errorf("Open error: %v\n", err)
			return
		}
		splitB := strings.Split(string(b), "\n")
		if len(splitB) == 1 {
			Line = append(Line, splitB[0])
		} else {
			Line = append(Line, splitB[0])
			joinLine := strings.Join(Line, "")
			fmt.Printf("read: %s\n", joinLine)
			Line = []string{}
			Line = append(Line, splitB[1])
		}
	}

}
