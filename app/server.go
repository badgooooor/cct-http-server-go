package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	req := make([]byte, 1024)
	conn.Read(req)
	data_str := string(req)

	request_line := strings.TrimSpace(data_str)
	parts := strings.Split(request_line, " ")

	// Path section
	if parts[1] != "/" {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		conn.Close()
	} else {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	}
}
