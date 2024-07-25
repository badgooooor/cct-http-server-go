package main

import (
	"fmt"
	"net"
	"net/http"
)

type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func (r HTTPResponse) String() string {
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, http.StatusText(r.StatusCode))
	for k, v := range r.Headers {
		response += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	response += "\r\n" + string(r.Body)
	return response
}

func (r HTTPResponse) Write(conn net.Conn) error {
	_, err := conn.Write([]byte(r.String()))
	return err
}

// Composition
func textResponse(statusCode int, body string) *HTTPResponse {
	response := &HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(body)),
		},
		Body: []byte(body),
	}

	return response
}

func fileResponse(statusCode int, body []byte) *HTTPResponse {
	response := &HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":   "application/octet-stream",
			"Content-Length": fmt.Sprintf("%d", len(body)),
		},
		Body: body,
	}

	return response
}
