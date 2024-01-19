package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:88")
	if err != nil {
		fmt.Printf("connection failed, err: %v\n", err)
	}
	defer conn.Close()

	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n')
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" {
			return
		}

		_, err := conn.Write([]byte(inputInfo))
		if err != nil {
			fmt.Printf("send failed, err: %v\n", err)
			return
		}
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("receive failed, err: %v\n", err)
			return
		}
		fmt.Printf("received: %v\n", string(buf[:n]))
	}
}