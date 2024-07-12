package main

import "strings"

type Request struct {
	data       string
	data_parts []string
}

func NewRequest(data string) *Request {
	request_line := strings.TrimSpace(data)
	data_parts := strings.Split(request_line, " ")

	return &Request{
		data,
		data_parts,
	}
}

func (r *Request) RawPath() string {
	return r.data_parts[1]
}

func (r *Request) Path() []string {
	return strings.Split(r.data_parts[1], "/")
}
