package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sidalikharef/goremote/drivers/registry"
	"github.com/sidalikharef/goremote/drivers/roku"
)

func TestRoutesBasic(t *testing.T) {
	r := registry.New()
	s := New(r)
	ts := httptest.NewServer(s.Routes())
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("bad status: %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestKeypressValidation(t *testing.T) {
	r := registry.New()
	s := New(r)
	ts := httptest.NewServer(s.Routes())
	defer ts.Close()
	resp, _ := http.Post(ts.URL+"/api/keypress/Home", "", nil)
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestKeypressWithMock(t *testing.T) {
	r := registry.New()
	d, _ := r.RegisterDevice("tv1", "TV", "roku", "192.168.1.100")
	c, _ := roku.New("192.168.1.100", time.Second)
	d.SetDriver(c)
	s := New(r)
	ts := httptest.NewServer(s.Routes())
	defer ts.Close()
	resp, _ := http.Post(ts.URL+"/api/keypress/BadKey", "", nil)
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 got %d", resp.StatusCode)
	}
	resp.Body.Close()
}
