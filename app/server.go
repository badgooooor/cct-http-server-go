package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Flags
	directoryPtr := flag.String("directory", "/tmp/", "file directory")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn, directoryPtr)
	}
}

func handleConnection(conn net.Conn, directoryPtr *string) {
	defer conn.Close()

	req := make([]byte, 1024)
	conn.Read(req)
	data_str := string(req)
	req_data := NewRequest(data_str)

	// Handle paths
	if req_data.RawPath() == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if req_data.Path()[1] == "echo" {
		message := req_data.Path()[2]
		response := textResponse(200, message)
		encoding := req_data.Headers()["Accept-Encoding"]
		if encoding == "gzip" {
			response.Headers["Content-Encoding"] = encoding
		}
		response.Write(conn)
	} else if req_data.Path()[1] == "user-agent" {
		message := req_data.Headers()["User-Agent"]
		response := textResponse(200, message)
		response.Write(conn)
	} else if req_data.Path()[1] == "files" {
		dir := *directoryPtr
		fileName := strings.TrimPrefix(req_data.Path()[2], "/files/")

		method := req_data.Method()

		if method == "GET" {
			data, err := os.ReadFile(dir + fileName)
			if err != nil {
				response := textResponse(404, "")
				response.Write(conn)
			} else {
				conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)))
			}
		} else if method == "POST" {
			file := []byte(strings.Trim(req_data.body, "\x00"))

			if err := os.WriteFile(dir+fileName, file, 0644); err != nil {
				response := textResponse(404, "")
				response.Write(conn)
			} else {
				conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
			}
		} else {
			response := textResponse(404, "")
			response.Write(conn)
		}
	} else {
		response := textResponse(404, "")
		response.Write(conn)
	}
}
