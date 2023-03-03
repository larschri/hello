package main

import (
	"context"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/chromedp/chromedp"
)

func TestIt(t *testing.T) {
	server := httptest.NewServer(helloHandler)
	defer server.Close()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(server.URL),
		chromedp.Text(`H1`, &res, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}

	if res != "Hello" {
		t.Errorf("Expected Hello, got %s", res)
	}
}
