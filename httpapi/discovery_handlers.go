package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sidalikharef/goremote/config"
	"github.com/sidalikharef/goremote/discovery"
	vidaaDisc "github.com/sidalikharef/goremote/discovery/vidaa"
	"github.com/sidalikharef/goremote/discovery/ssdp"
	"github.com/sidalikharef/goremote/drivers/roku"
)

var discManager = discovery.NewManager()

func (s *Server) handleDiscoverStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ssdpD := ssdp.New()
	vidaaP := vidaaDisc.New()
	orch := discovery.NewOrchestrator(discManager, ssdpD, vidaaP)
	ctx := r.Context()
	res, errs := orch.Discover(ctx, discovery.Options{Timeout: 4 * time.Second, MaxResults: 50})
	// If providers error but return partial results, still return what we have
	if len(res) == 0 && len(errs) > 0 {
		// Still return empty array rather than 500 to not block UI
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"devices": []interface{}{}, "errors": len(errs)})
		return
	}
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
	existing, err := config.Load(config.DefaultPath())
	if err != nil {
		existing = &config.Config{Version: 1, Devices: []config.Device{}}
	}
	// replace or add
	found := false
	for i := range existing.Devices {
		if existing.Devices[i].ID == req.ID {
			existing.Devices[i].Name = req.Name
			existing.Devices[i].Driver = req.Driver
			existing.Devices[i].Address = req.Address
			found = true
			break
		}
	}
	if !found {
		existing.Devices = append(existing.Devices, config.Device{ID: req.ID, Name: req.Name, Driver: req.Driver, Address: req.Address})
	}
	if existing.Active == "" {
		existing.Active = req.ID
	}
	_ = config.Save(config.DefaultPath(), existing)
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
