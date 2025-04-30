package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	upAdd, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Resolve UDP Addr: %v\n", err)
		return
	}
	upCon, err := net.DialUDP("udp", nil, upAdd)
	if err != nil {
		fmt.Printf("DialUDP: %v\n", err)
		return
	}
	defer upCon.Close()
	stdinReader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print(">")
		strRead, err := stdinReader.ReadString('\n')
		if err != nil {
			if strRead != "" {
				_, err = upCon.Write([]byte(strRead))
				if err != nil {
					fmt.Printf("upCon.Write: %v\n", err)
				}
			}
			fmt.Printf("Read error: %v\n", err)
			return
		}
		_, err = upCon.Write([]byte(strRead))
		if err != nil {
			fmt.Printf("upCon.Write: %v\n", err)
		}
	}

}
