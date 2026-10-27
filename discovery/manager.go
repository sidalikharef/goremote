package discovery

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/sidalikharef/goremote/config"
	"github.com/sidalikharef/goremote/internal/validate"
)

// Manager coordinates discovery.
type Manager struct {
	mu   sync.RWMutex
	seen map[string]Result // keyed by stable ID
}

func NewManager() *Manager {
	return &Manager{
		seen: map[string]Result{},
	}
}

// Deduplicate by stable ID.
func (m *Manager) Deduplicate(results []Result) []Result {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []Result{}
	for _, r := range results {
		if !Validate(r) {
			continue
		}
		if _, exists := m.seen[r.ID]; exists {
			continue
		}
		m.seen[r.ID] = r
		out = append(out, r)
	}
	return out
}

// HandleIPChange checks if IP changed.
func (m *Manager) HandleIPChange(oldID, newAddr string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if r, ok := m.seen[oldID]; ok {
		return r.Address != newAddr
	}
	return false
}

// AddManual adds a manually entered device (fallback).
func (m *Manager) AddManual(id, name, driver, addr string) (Result, error) {
	if net.ParseIP(addr) == nil {
		return Result{}, fmt.Errorf("invalid address")
	}
	if err := validate.IP(addr); err != nil {
		if net.ParseIP(addr) == nil {
			return Result{}, fmt.Errorf("invalid address")
		}
	}
	if id == "" {
		id = addr
	}
	r := Result{
		ID:         id,
		Name:       name,
		DriverHint: driver,
		Address:    addr,
		Source:     "manual",
	}
	if !Validate(r) {
		return Result{}, fmt.Errorf("invalid device")
	}
	m.mu.Lock()
	m.seen[r.ID] = r
	m.mu.Unlock()
	return r, nil
}

// ToConfigDevice converts to config device - requires explicit user action.
func ToConfigDevice(r Result) (config.Device, error) {
	if !Validate(r) {
		return config.Device{}, fmt.Errorf("invalid")
	}
	if net.ParseIP(r.Address) == nil {
		return config.Device{}, fmt.Errorf("invalid address")
	}
	return config.Device{
		ID:      r.ID,
		Name:    r.Name,
		Driver:  r.DriverHint,
		Address: r.Address,
		Added:   time.Now(),
	}, nil
}
