package discovery

import "context"

// Provider is a discovery provider interface.
type Provider interface {
	Name() string
	Discover(ctx context.Context, opts Options) ([]Result, error)
}
