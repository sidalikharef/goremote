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
	ID         string
	Name       string
	DriverHint string
	Address    string
	Location   string
	USN        string
	Model      string
	Firmware   string
	Source     string // ssdp, mdns, manual
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
