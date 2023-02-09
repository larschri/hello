package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -generate types,chi-server,spec -package main -o generated.go hello.yaml

type HelloImpl struct{}

func (*HelloImpl) GetHelloName(w http.ResponseWriter, r *http.Request, name string) {
	json.NewEncoder(w).Encode(Message{"Hello, " + name + "!"})
}

func main() {
	var myApi HelloImpl

	ts := httptest.NewServer(Handler(&myApi))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/hello/oapi")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}
