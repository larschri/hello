package main

import (
	"net"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
}

// initialize sets up listeners and handlers
func initialize(addr string) net.Listener {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/foo", http.HandlerFunc(handler))
	http.HandleFunc("/bar", http.HandlerFunc(handler))

	return l
}

func main() {
	http.Serve(initialize(":8080"), nil)
}
