package main

import (
	"fmt"
	"strings"
)

type Headers map[string]string

type Request struct {
	data_parts []string
	headers    Headers
}

func NewRequest(data string) *Request {
	fmt.Printf("%s\n", data)
	request_line := strings.TrimSpace(data)
	data_parts := strings.Split(request_line, " ")

	headers := parseHeaders(data_parts)
	for k, v := range headers {
		fmt.Println(k, v)
	}
	return &Request{
		data_parts,
		headers,
	}
}

func parseHeaders(data_parts []string) Headers {
	headers := make(map[string]string)
	for i := 2; i < len(data_parts)-1; i++ {
		splited_data_key := strings.Split(data_parts[i], "\n")
		key := strings.Replace(splited_data_key[1], ":", "", 1)
		splited_data_key_value := strings.Split(data_parts[i+1], "\n")
		value := strings.TrimSpace(splited_data_key_value[0])
		headers[key] = value
	}

	return headers
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
