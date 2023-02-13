package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var atomicmap = atomic.Pointer[map[string]string]{}

func updateMap() {
	mp := make(map[string]string)
	mp["hello"] = time.Now().String()
	mp["tomorrow"] = time.Now().Add(24 * time.Hour).String()
	atomicmap.Store(&mp)
}

func readMap() string {
	return fmt.Sprintf("%v", *atomicmap.Load())
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for ctx.Err() == nil {
				updateMap()
				readMap()
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
