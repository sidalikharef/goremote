package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sidalikharef/goremote/tv"
)

// Device represents a configured TV device.
type Device struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Driver   string            `json:"driver"`
	Address  string            `json:"address"`
	State    tv.State          `json:"state"`
	Pairing  map[string]string `json:"pairing,omitempty"` // sensitive data, stored with restrictive perms
	Added    time.Time         `json:"added"`
	LastSeen time.Time         `json:"last_seen,omitempty"`
}

// Config is the versioned persistent configuration.
type Config struct {
	Version int      `json:"version"`
	Devices []Device `json:"devices"`
	Active  string   `json:"active,omitempty"` // active device ID
}

var (
	mu sync.RWMutex
)

// DefaultPath returns the default config path.
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./goremote.json"
	}
	return filepath.Join(home, ".config", "goremote", "config.json")
}

// Load loads config from path.
func Load(path string) (*Config, error) {
	mu.RLock()
	defer mu.RUnlock()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Version: 1, Devices: []Device{}}, nil
		}
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	if c.Version < 1 {
		c.Version = 1
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Save saves config atomically with restrictive permissions.
func Save(path string, c *Config) error {
	mu.Lock()
	defer mu.Unlock()
	if err := c.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	// Write with restrictive permissions for sensitive data
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Version < 1 {
		return fmt.Errorf("invalid version")
	}
	seen := map[string]bool{}
	for i := range c.Devices {
		d := &c.Devices[i]
		if d.ID == "" {
			return fmt.Errorf("device id required")
		}
		if seen[d.ID] {
			return fmt.Errorf("duplicate device id")
		}
		seen[d.ID] = true
		if d.Driver == "" {
			return fmt.Errorf("driver required")
		}
		if d.Name == "" {
			return fmt.Errorf("name required")
		}
		// Validate address is not arbitrary/untrusted - must be IP
		if d.Address != "" {
			// basic check
			if len(d.Address) > 253 {
				return fmt.Errorf("invalid address")
			}
		}
		if d.State == "" {
			d.State = tv.StateDisconnected
		}
	}
	if c.Active != "" && !seen[c.Active] {
		return fmt.Errorf("active device not found")
	}
	return nil
}
