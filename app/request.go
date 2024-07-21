package main

import (
	"slices"
	"strings"
)

type Headers map[string]string

type Request struct {
	data_parts []string
	headers    Headers
	body       string
}

func NewRequest(data string) *Request {
	request_line := strings.TrimSpace(data)
	data_parts := strings.Split(request_line, " ")

	body, seperator_idx := parseBody(request_line)
	headers := parseHeaders(data_parts, seperator_idx)
	return &Request{
		data_parts,
		headers,
		body,
	}
}

func parseHeaders(data_parts []string, seperator_idx int) Headers {
	headers := make(map[string]string)
	for i := 2; i < seperator_idx+1; i++ {
		splited_data_key := strings.Split(data_parts[i], "\n")

		key := strings.Replace(splited_data_key[1], ":", "", 1)

		splited_data_key_value := strings.Split(data_parts[i+1], "\n")
		value := strings.TrimSpace(splited_data_key_value[0])

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
	return r.data_parts[0]
}

func (r *Request) RawPath() string {
	return r.data_parts[1]
}

func (r *Request) Path() []string {
	return strings.Split(r.data_parts[1], "/")
}

func (r *Request) Headers() Headers {
	return r.headers
}
