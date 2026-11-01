package pocketsmith

import (
	"context"
	"net/http"
	"testing"
)

func TestProfileResolvesTheTimeZone(t *testing.T) {
	c, _ := server(t, map[string]http.HandlerFunc{
		"GET /me": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"id":42,"time_zone":"Eastern Time (US & Canada)"}`))
		},
		"GET /time_zones": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`[{"name":"Hawaii","identifier":"Pacific/Honolulu"},{"name":"Eastern Time (US & Canada)","identifier":"America/New_York"}]`))
		},
	})
	p, err := c.Profile(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != 42 || p.TimeZone.String() != "America/New_York" {
		t.Errorf("Profile = %d %s, want 42 America/New_York", p.ID, p.TimeZone)
	}
}

func TestProfileFailsOnAnUnknownZone(t *testing.T) {
	c, _ := server(t, map[string]http.HandlerFunc{
		"GET /me":         func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"id":1,"time_zone":"Nowhere"}`)) },
		"GET /time_zones": func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`[]`)) },
	})
	if _, err := c.Profile(context.Background()); err == nil {
		t.Fatal("Profile succeeded for an unknown zone")
	}
}
