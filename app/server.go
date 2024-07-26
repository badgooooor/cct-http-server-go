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

	router := NewRouter()

	router.Handle("GET", "/", func(r Request) Response {
		return Response{StatusCode: 200}
	})

	router.Handle("GET", "/echo/*", func(r Request) Response {
		var response *Response
		message := r.Path()[2]

		encoding := r.Headers()["Accept-Encoding"]
		if strings.Contains(encoding, "gzip") {
			response = textResponseWithOpts(200, message, &TextResponseOpts{
				Compression: "gzip",
			})
			response.Headers["Content-Encoding"] = "gzip"
		} else {
			response = textResponse(200, message)
		}

		return *response
	})

	router.Handle("GET", "/user-agent", func(r Request) Response {
		message := r.Headers()["User-Agent"]
		response := textResponse(200, message)
		return *response
	})

	router.Handle("GET", "/files/*", func(r Request) Response {
		dir := *directoryPtr
		fileName := strings.TrimPrefix(r.Path()[2], "/files/")

		data, err := os.ReadFile(dir + fileName)
		var response *Response
		if err != nil {
			response = textResponse(404, "")
		} else {
			response = fileResponse(200, data)
		}

		return *response
	})

	router.Handle("POST", "/files/*", func(r Request) Response {
		dir := *directoryPtr
		fileName := strings.TrimPrefix(r.Path()[2], "/files/")

		file := []byte(strings.Trim(r.body, "\x00"))

		var response *Response
		if err := os.WriteFile(dir+fileName, file, 0644); err != nil {
			response = textResponse(404, "")

		} else {
			response = textResponse(201, "")
		}

		return *response
	})

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go router.handleConnection(conn)
	}
}
