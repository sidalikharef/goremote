package roku

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRokuKeyWithMockServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			return
		}
		if !strings.Contains(r.URL.Path, "/keypress/Home") {
			w.WriteHeader(400)
			return
		}
		w.WriteHeader(200)
	}))
	defer ts.Close()
	// Can't easily override host/port without modifying client; test validation
	c, err := New("192.168.1.100", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if c.Name() != "roku" {
		t.Fatal("bad name")
	}
	if err := c.Key(context.Background(), "BadKey!!!"); err == nil {
		t.Fatal("should fail")
	}
	if err := c.Launch(context.Background(), "../bad"); err == nil {
		t.Fatal("should fail")
	}
	caps := c.Capabilities()
	if !caps.Navigation {
		t.Fatal("bad caps")
	}
}
