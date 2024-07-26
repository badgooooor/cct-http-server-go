package main

import (
	"fmt"
	"net"
	"strings"
)

type Handler func(Request) Response

type Router struct {
	routes map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		routes: map[string]Handler{},
	}
}

func (r *Router) Handle(method string, endpoint string, handler Handler) {
	r.routes[method+endpoint] = handler
}

func (r *Router) GetHandler(method string, endpoint string) Handler {
	if handler, ok := r.routes[method+endpoint]; ok {
		return handler
	}

	for routeInfo := range r.routes {
		if strings.HasSuffix(routeInfo, "*") && strings.HasPrefix(method+endpoint, routeInfo[:len(routeInfo)-1]) {
			return r.routes[routeInfo]
		}
	}
	return nil
}

func (r *Router) handleConnection(conn net.Conn) error {
	defer conn.Close()

	req := make([]byte, 1024)
	conn.Read(req)
	dataString := string(req)
	request := NewRequest(dataString)

	handler := r.GetHandler(request.Method(), request.RawPath())
	var response Response
	if handler == nil {
		response = Response{StatusCode: 404}
	} else {
		response = handler(*request)
	}
	if err := response.Write(conn); err != nil {
		fmt.Println("Error writing response: ", err.Error())
		return err
	}
	fmt.Println("Processed request: ", request.Method(), request.RawPath())
	return nil
}
