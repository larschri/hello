package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"
)

//go:embed static
var embeddedFS embed.FS

func triggerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Trigger", `{"myTrigger":"value passed in response header"}`)
	fmt.Fprint(w, "response body from /trigger")
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Println(k, v)
		if !strings.HasPrefix(strings.ToUpper(k), "HX-") {
			continue
		}
		fmt.Fprintf(w, "%v: %#v\n", k, v)
	}
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprint(w, `This response took 3 seconds`)

}

func main() {
	serverRoot, err := fs.Sub(embeddedFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/headers", headersHandler)
	http.HandleFunc("/trigger", triggerHandler)
	http.HandleFunc("/slow", slowHandler)
	http.Handle("/", http.FileServer(http.FS(serverRoot)))
	fmt.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
