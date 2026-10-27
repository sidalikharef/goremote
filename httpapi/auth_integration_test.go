package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sidalikharef/goremote/drivers/registry"
	"github.com/sidalikharef/goremote/internal/auth"
	"github.com/sidalikharef/goremote/tv"
)

func setupTest(t *testing.T) (*auth.Store, *registry.Registry, http.Handler) {
	dir := t.TempDir()
	as := auth.New(dir + "/auth.json")
	r := registry.New()
	s := NewWithAuth(r, as)
	h := SecurityHeaders(s.Routes())
	h = MaxBody(1024 * 64)(h)
	allowed := map[string]bool{"localhost:8080": true, "127.0.0.1:8080": true}
	h = CheckHostSafe(allowed)(h)
	h = CheckOriginSafe(allowed)(h)
	h = RequireAuth(as)(h)
	return as, r, h
}

func TestFreshInstallStatus(t *testing.T) {
	_, _, h := setupTest(t)
	ts := httptest.NewServer(h)
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/auth/status", nil)
	req.Host = "localhost:8080"
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestUnauthorizedCommand(t *testing.T) {
	_, _, h := setupTest(t)
	ts := httptest.NewServer(h)
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/keypress/Home", nil)
	req.Host = "localhost:8080"
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 401 {
		t.Fatalf("want 401 got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestEnrollFlow(t *testing.T) {
	as, _, h := setupTest(t)
	ts := httptest.NewServer(h)
	defer ts.Close()
	token, _ := as.GenerateBootstrapToken()
	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "longenoughpass", "bootstrap_token": token})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/auth/enroll", bytes.NewReader(body))
	req.Host = "127.0.0.1:8080"
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 201 {
		t.Fatalf("want 201 got %d", resp.StatusCode)
	}
	resp.Body.Close()
	body2, _ := json.Marshal(map[string]string{"username": "admin2", "password": "longenoughpass", "bootstrap_token": token})
	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/auth/enroll", bytes.NewReader(body2))
	req2.Host = "127.0.0.1:8080"
	req2.RemoteAddr = "127.0.0.1:12345"
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := http.DefaultClient.Do(req2)
	if resp2.StatusCode == 201 {
		t.Fatal("should reject second enroll")
	}
	resp2.Body.Close()
}

func TestLoginLogout(t *testing.T) {
	as, _, h := setupTest(t)
	ts := httptest.NewServer(h)
	defer ts.Close()
	token, _ := as.GenerateBootstrapToken()
	as.Enroll("admin", "longenoughpass", token)
	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "longenoughpass"})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/auth/login", bytes.NewReader(body))
	req.Host = "localhost:8080"
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp.StatusCode)
	}
	var sessionCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "goremote_session" {
			sessionCookie = c
		}
	}
	if sessionCookie == nil {
		t.Fatal("no cookie")
	}
	resp.Body.Close()
	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/auth/logout", nil)
	req2.Host = "localhost:8080"
	req2.AddCookie(sessionCookie)
	resp2, _ := http.DefaultClient.Do(req2)
	if resp2.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp2.StatusCode)
	}
	resp2.Body.Close()
}

func TestWrongPassword(t *testing.T) {
	as, _, h := setupTest(t)
	ts := httptest.NewServer(h)
	defer ts.Close()
	token, _ := as.GenerateBootstrapToken()
	as.Enroll("admin", "longenoughpass", token)
	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "wrongpass"})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/auth/login", bytes.NewReader(body))
	req.Host = "localhost:8080"
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 401 {
		t.Fatalf("want 401 got %d", resp.StatusCode)
	}
	resp.Body.Close()
}

type mockDrv struct {
	keyCalled bool
}

func (m *mockDrv) Name() string                                { return "mock" }
func (m *mockDrv) Key(ctx context.Context, k string) error     { m.keyCalled = true; return nil }
func (m *mockDrv) Apps(ctx context.Context) ([]tv.App, error)  { return []tv.App{}, nil }
func (m *mockDrv) Launch(ctx context.Context, id string) error { return nil }

func TestAuthorizedMockRoku(t *testing.T) {
	as, r, h := setupTest(t)
	ts := httptest.NewServer(h)
	defer ts.Close()
	token, _ := as.GenerateBootstrapToken()
	as.Enroll("admin", "longenoughpass", token)
	d, _ := r.RegisterDevice("tv1", "TV", "roku", "192.168.1.100")
	md := &mockDrv{}
	d.SetDriver(md)
	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "longenoughpass"})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/auth/login", bytes.NewReader(body))
	req.Host = "localhost:8080"
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var sessionCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "goremote_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("no cookie")
	}
	resp.Body.Close()
	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/keypress/Home?device=tv1", nil)
	req2.Host = "localhost:8080"
	req2.AddCookie(sessionCookie)
	resp2, _ := http.DefaultClient.Do(req2)
	if resp2.StatusCode != 200 {
		t.Fatalf("want 200 got %d", resp2.StatusCode)
	}
	resp2.Body.Close()
	if !md.keyCalled {
		t.Fatal("key not called")
	}
}

func TestInvalidDeviceFallback(t *testing.T) {
	as, r, h := setupTest(t)
	ts := httptest.NewServer(h)
	defer ts.Close()
	token, _ := as.GenerateBootstrapToken()
	as.Enroll("admin", "longenoughpass", token)
	d1, _ := r.RegisterDevice("tv1", "TV1", "roku", "192.168.1.100")
	d1.SetDriver(&mockDrv{})
	d2, _ := r.RegisterDevice("tv2", "TV2", "roku", "192.168.1.101")
	md2 := &mockDrv{}
	d2.SetDriver(md2)
	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "longenoughpass"})
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/auth/login", bytes.NewReader(body))
	req.Host = "localhost:8080"
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var sessionCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "goremote_session" {
			sessionCookie = c
			break
		}
	}
	resp.Body.Close()
	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/keypress/Home?device=badid", nil)
	req2.Host = "localhost:8080"
	req2.AddCookie(sessionCookie)
	resp3, _ := http.DefaultClient.Do(req2)
	if resp3.StatusCode == 200 {
		t.Fatal("should not succeed with invalid device")
	}
	resp3.Body.Close()
	if md2.keyCalled {
		t.Fatal("should not call wrong device")
	}
}
