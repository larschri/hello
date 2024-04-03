package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	go http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintln(w, "hello world")
	}))

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			resp, err := http.Get("http://localhost:8080/")
			if err != nil {
				panic(err)
			}
			if _, err := io.ReadAll(resp.Body); err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
