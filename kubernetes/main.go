package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Kubernetes")
}

func main() {
	log.Println("Hello, kubernetes")
	http.ListenAndServe(":8080", http.HandlerFunc(hello))
}
