package pocketsmith

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const key = "test-key"

// server serves handlers, checking every request carries the developer key, and returns a
// client pointed at it.
func server(t *testing.T, handlers map[string]http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	mux := http.NewServeMux()
	for pattern, h := range handlers {
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Developer-Key") != key {
				t.Errorf("%s %s without the developer key", r.Method, r.URL.Path)
			}
			h(w, r)
		})
	}
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return &Client{BaseURL: srv.URL, Key: key, HTTP: srv.Client()}, srv
}

func TestErrorStatusIsReportedWithoutTheKey(t *testing.T) {
	c, _ := server(t, map[string]http.HandlerFunc{
		"GET /me": func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "<html>bad gateway</html>", http.StatusBadGateway)
		},
	})
	_, err := c.Profile(context.Background())
	if err == nil || !strings.Contains(err.Error(), "502") || !strings.Contains(err.Error(), "GET /me") {
		t.Fatalf("err = %v, want GET /me and a 502", err)
	}
	if strings.Contains(err.Error(), key) {
		t.Errorf("error leaks the key: %v", err)
	}
}
