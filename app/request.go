package main

import (
	"slices"
	"strings"
)

type Headers map[string]string

type Request struct {
	dataParts []string
	headers   Headers
	body      string
}

func NewRequest(data string) *Request {
	requestLine := strings.TrimSpace(data)
	body, seperator_idx := parseBody(requestLine)

	dataParts := strings.Split(requestLine, "\r\n")
	headers := parseHeaders(dataParts, seperator_idx)

	return &Request{
		dataParts,
		headers,
		body,
	}
}

func parseHeaders(dataParts []string, seperatorIdx int) Headers {
	headers := make(map[string]string)
	for i := 2; i < seperatorIdx; i++ {
		splitedLine := strings.Split(dataParts[i], ":")

		key := splitedLine[0]
		value := strings.TrimSpace(splitedLine[1])
		headers[key] = value
	}

	return headers
}

// Extract request body from HTTP request with seperator line index for extract headers.
func parseBody(request_line string) (string, int) {
	lines := strings.Split(request_line, "\r\n")
	seperator_idx := slices.IndexFunc(lines, func(item string) bool {
		return len(item) == 0
	})

	return strings.Join(lines[seperator_idx+1:], ""), seperator_idx
}

func (r *Request) Method() string {
	splited := strings.Split(r.dataParts[0], " ")
	return splited[0]
}

func (r *Request) RawPath() string {
	splited := strings.Split(r.dataParts[0], " ")
	return splited[1]
}

func (r *Request) Path() []string {
	return strings.Split(r.RawPath(), "/")
}

func (r *Request) Headers() Headers {
	return r.headers
}
