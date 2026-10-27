package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sidalikharef/goremote/config"
	"github.com/sidalikharef/goremote/discovery"
	"github.com/sidalikharef/goremote/discovery/ssdp"
	"github.com/sidalikharef/goremote/drivers/roku"
)

var discManager = discovery.NewManager()

func (s *Server) handleDiscoverStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	d := ssdp.New()
	ctx := r.Context()
	res, err := d.DiscoverRoku(ctx, discovery.Options{Timeout: 3 * time.Second, MaxResults: 20})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res = discManager.Deduplicate(res)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"devices": res})
}

func (s *Server) handleDiscoverList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"devices": []interface{}{}})
}

func (s *Server) handleDiscoverSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID      string `json:"id"`
		Address string `json:"address"`
		Name    string `json:"name"`
		Driver  string `json:"driver"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Address == "" || req.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := s.reg.RegisterDevice(req.ID, req.Name, req.Driver, req.Address); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Driver == "roku" || req.Driver == "" {
		if c, err := roku.New(req.Address, 5*time.Second); err == nil {
			if dev, ok := s.reg.GetDevice(req.ID); ok {
				dev.SetDriver(c)
			}
		}
	}
	cfg := &config.Config{Version: 1, Devices: []config.Device{{ID: req.ID, Name: req.Name, Driver: "roku", Address: req.Address}}, Active: req.ID}
	_ = config.Save(config.DefaultPath(), cfg)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleDevicesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	devs := s.reg.ListDevices()
	type resp struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Driver string `json:"driver"`
		State  string `json:"state"`
	}
	res := make([]resp, 0, len(devs))
	for _, d := range devs {
		res = append(res, resp{ID: d.ID, Name: d.Name, Driver: d.Driver, State: string(d.GetState())})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
