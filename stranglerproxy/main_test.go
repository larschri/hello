package main

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"testing"
)

// stranglerProxy wraps a legacy service.
type stranglerProxy struct {
	*httputil.ReverseProxy
}

var strangler stranglerProxy

// ServeHTTP is used to strangle an old legacy service piece by piece. Some
// requests are sent to the legacy service and other requests are handled by
// new logic that replaces the legacy service.
func (s stranglerProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/bar" {
		w.Write([]byte("Hello, strangler"))
		return
	}

	s.ReverseProxy.ServeHTTP(w, r)
}

func TestMain(m *testing.M) {
	l := initialize(":0")
	defer l.Close()

	go http.Serve(l, nil)
	strangler.ReverseProxy = httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   l.Addr().String(),
	})

	os.Exit(m.Run())
}

func TestStranglerProxy(t *testing.T) {
	tests := map[string]string{
		"/bar": "Hello, strangler",
		"/foo": "Hello, world",
	}

	for u, b := range tests {
		w := httptest.NewRecorder()
		strangler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, u, nil))
		if w.Body.String() != b {
			t.Errorf("unexpected body for %s: %s", u, w.Body.String())
		}
	}
}
