package discovery

import (
	"context"
	"net"
	"time"

	"github.com/sidalikharef/goremote/tv"
)

// Options configures discovery.
type Options struct {
	Timeout      time.Duration
	MaxResults   int
	NetworkScope string // e.g., "lan"
}

// Result represents a discovered device.
type Result struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	DriverHint     string `json:"driver_hint,omitempty"`
	Driver         string `json:"driver,omitempty"`
	Manufacturer   string `json:"manufacturer,omitempty"`
	Model          string `json:"model,omitempty"`
	Firmware       string `json:"firmware,omitempty"`
	Address        string `json:"address"`
	Location       string `json:"location,omitempty"`
	USN            string `json:"usn,omitempty"`
	Source         string `json:"discovery_source,omitempty"`
	DiscoverySource string `json:"discovery_source_alias,omitempty"`
	Confidence     string `json:"confidence,omitempty"`
	Experimental   bool   `json:"experimental,omitempty"`
	RequiresPairing bool  `json:"requires_pairing,omitempty"`
}

// Service defines discovery service interface.
type Service interface {
	Discover(ctx context.Context, opts Options) ([]Result, error)
}

// Store tracks discovered devices for deduplication.
type Store struct {
}

// Validate checks if a device is a valid TV candidate.
func Validate(r Result) bool {
	if r.Address == "" {
		return false
	}
	if net.ParseIP(r.Address) == nil {
		return false
	}
	if r.ID == "" && r.USN == "" {
		return false
	}
	return true
}

// ToDiscovered converts to tv.DiscoveredDevice.
func ToDiscovered(r Result) tv.DiscoveredDevice {
	return tv.DiscoveredDevice{
		ID:       r.ID,
		Name:     r.Name,
		Driver:   r.DriverHint,
		Address:  r.Address,
		Model:    r.Model,
		Firmware: r.Firmware,
	}
}
