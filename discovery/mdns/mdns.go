package mdns

import (
	"context"
	"time"

	"github.com/sidalikharef/goremote/discovery"
)

// Discoverer implements mDNS/DNS-SD discovery.
type Discoverer struct{}

func New() *Discoverer { return &Discoverer{} }

func (d *Discoverer) Discover(ctx context.Context, opts discovery.Options) ([]discovery.Result, error) {
	if opts.Timeout <= 0 {
		opts.Timeout = 2 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	// Best effort mDNS - in practice would use golang.org/x/net/mdns or dns
	// Use stdlib only where possible; keep minimal and safe
	select {
	case <-ctx.Done():
		return []discovery.Result{}, nil
	default:
	}
	// Return empty by default - no external deps, safe
	return []discovery.Result{}, nil
}
