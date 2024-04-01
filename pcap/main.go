package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	server := &http.Server{Addr: os.Args[1], Handler: nil}
	http.Handle("/exit", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go server.Shutdown(context.Background())
	}))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
		if f, ok := w.(http.Flusher); ok {
			// flush to send a tcp packet already
			f.Flush()
			// sleep to allow another request send packets between ours
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Fprintln(w, "world")
	}))
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
