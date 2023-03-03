package main

import (
	"fmt"
	"net/http"
)

var helloHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<html><body><h1>Hello</h1>lorem ipsum</body></html>`)
}

func main() {
	http.ListenAndServe(":9042", helloHandler)
}
