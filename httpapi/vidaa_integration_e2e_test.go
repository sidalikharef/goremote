package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sidalikharef/goremote/drivers/registry"
	vidaa "github.com/sidalikharef/goremote/drivers/vidaa"
	"github.com/sidalikharef/goremote/internal/auth"
)

func TestVidaaHTTPPathRegistered(t *testing.T) {
	vidaa.EnableExperimental()
	defer vidaa.DisableExperimental()
	r := registry.New()
	r.RegisterFactory("vidaa", vidaa.Factory)
	as := auth.New(t.TempDir() + "/auth.json")
	s := NewWithAuth(r, as)
	h := SecurityHeaders(s.Routes())
	ts := httptest.NewServer(h)
	defer ts.Close()
	resp, _ := http.Post(ts.URL+"/api/devices/vidaa/pair/start", "application/json", nil)
	if resp.StatusCode == 200 {
		t.Fatal("expected auth required")
	}
}
