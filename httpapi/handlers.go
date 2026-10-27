package httpapi

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/sidalikharef/goremote/drivers/registry"
	"github.com/sidalikharef/goremote/internal/auth"
	"github.com/sidalikharef/goremote/internal/validate"
	"github.com/sidalikharef/goremote/tv"
)

//go:embed all:web
var webFS embed.FS

// Server handles HTTP requests.
type Server struct {
	reg  *registry.Registry
	auth *auth.Store
}

func New(reg *registry.Registry) *Server {
	return &Server{reg: reg}
}

func NewWithAuth(reg *registry.Registry, a *auth.Store) *Server {
	return &Server{reg: reg, auth: a}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	sub, _ := fs.Sub(webFS, "web")
	mux.Handle("/", http.FileServer(http.FS(sub)))
	if s.auth != nil {
		as := NewAuth(s.auth)
		MountAuth(mux, as)
	}
	mux.HandleFunc("/api/keypress/", s.handleKeypress)
	mux.HandleFunc("/api/apps", s.handleApps)
	mux.HandleFunc("/api/launch/", s.handleLaunch)
	mux.HandleFunc("/api/devices", s.handleDevicesList)
	mux.HandleFunc("/api/discover", s.handleDiscoverStart)
	mux.HandleFunc("/api/discover/start", s.handleDiscoverStart)
	mux.HandleFunc("/api/discover/list", s.handleDiscoverList)
	mux.HandleFunc("/api/discover/save", s.handleDiscoverSave)
	return mux
}

func (s *Server) getActiveDevice(r *http.Request) *registry.Device {
	id := r.URL.Query().Get("device")
	if id != "" {
		if d, ok := s.reg.GetDevice(id); ok {
			return d
		}
		return nil
	}
	devs := s.reg.ListDevices()
	if len(devs) > 0 {
		return devs[0]
	}
	return nil
}

func (s *Server) handleKeypress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	if strings.Contains(path, "..") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	key := strings.TrimPrefix(path, "/api/keypress/")
	if err := validate.KeyName(key); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	d := s.getActiveDevice(r)
	if d == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := d.Call(ctx, func(drv tv.Driver) error {
		return drv.Key(ctx, key)
	}); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleApps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	d := s.getActiveDevice(r)
	if d == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var apps []tv.App
	err := d.Call(ctx, func(drv tv.Driver) error {
		a, e := drv.Apps(ctx)
		if e != nil {
			return e
		}
		apps = a
		return nil
	})
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

func (s *Server) handleLaunch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if strings.Contains(r.URL.Path, "..") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/launch/")
	if err := validate.AppID(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	d := s.getActiveDevice(r)
	if d == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := d.Call(ctx, func(drv tv.Driver) error {
		return drv.Launch(ctx, id)
	}); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusOK)
}
