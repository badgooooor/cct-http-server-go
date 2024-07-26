package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"net/http"
)

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func (r Response) String() string {
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, http.StatusText(r.StatusCode))
	for k, v := range r.Headers {
		response += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	response += "\r\n" + string(r.Body)
	return response
}

func (r Response) Write(conn net.Conn) error {
	_, err := conn.Write([]byte(r.String()))
	return err
}

type TextResponseOpts struct {
	Compression string
}

func textResponse(statusCode int, body string) *Response {
	return textResponseWithOpts(statusCode, body, &TextResponseOpts{})
}

func textResponseWithOpts(statusCode int, body string, opts *TextResponseOpts) *Response {
	var content string
	if opts.Compression == "gzip" {
		var buffer bytes.Buffer
		w := gzip.NewWriter(&buffer)
		w.Write([]byte(body))
		w.Close()

		content = buffer.String()
	} else {
		content = body
	}

	response := &Response{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(content)),
		},
		Body: []byte(content),
	}

	return response
}

func fileResponse(statusCode int, body []byte) *Response {
	response := &Response{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":   "application/octet-stream",
			"Content-Length": fmt.Sprintf("%d", len(body)),
		},
		Body: body,
	}

	return response
}
