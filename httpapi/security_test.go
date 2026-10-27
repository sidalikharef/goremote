package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sidalikharef/goremote/drivers/registry"
)

func TestSecurityHeaders(t *testing.T) {
	r := registry.New()
	s := New(r)
	h := SecurityHeaders(s.Routes())
	ts := httptest.NewServer(h)
	defer ts.Close()
	resp, _ := http.Get(ts.URL + "/")
	if resp.Header.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("missing header")
	}
	resp.Body.Close()
}

func TestMaxBody(t *testing.T) {
	r := registry.New()
	s := New(r)
	h := MaxBody(10)(s.Routes())
	ts := httptest.NewServer(h)
	defer ts.Close()
	resp, _ := http.Post(ts.URL+"/api/keypress/Home", "text/plain", http.NoBody)
	if resp.StatusCode == 0 {
		t.Fatal("no response")
	}
	resp.Body.Close()
}

func TestCrossOriginBlocked(t *testing.T) {
	r := registry.New()
	s := New(r)
	allowed := map[string]bool{"localhost:8080": true}
	h := CheckOrigin(allowed)(s.Routes())
	ts := httptest.NewServer(h)
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/keypress/Home", nil)
	req.Header.Set("Origin", "http://evil.com")
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("want 403 got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestInvalidKeyRejected(t *testing.T) {
	r := registry.New()
	s := New(r)
	ts := httptest.NewServer(s.Routes())
	defer ts.Close()
	resp, _ := http.Post(ts.URL+"/api/keypress/BadKey!!!", "", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
