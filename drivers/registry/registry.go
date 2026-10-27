package registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/sidalikharef/goremote/tv"
)

// Factory creates a driver instance.
type Factory func(opts map[string]string) (tv.Driver, error)

// Registry manages drivers and devices.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
	devices   map[string]*Device
}

// Device represents a registered TV device instance.
type Device struct {
	mu     sync.RWMutex
	ID     string
	Name   string
	Driver string
	Addr   string
	State  tv.State
	Drv    tv.Driver
}

// New creates a new registry.
func New() *Registry {
	return &Registry{
		factories: map[string]Factory{},
		devices:   map[string]*Device{},
	}
}

// RegisterFactory registers a driver factory.
func (r *Registry) RegisterFactory(name string, f Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = f
}

// RegisterDevice registers a device.
func (r *Registry) RegisterDevice(id, name, driver, addr string) (*Device, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.devices[id]; ok {
		return nil, fmt.Errorf("device exists")
	}
	d := &Device{
		ID:     id,
		Name:   name,
		Driver: driver,
		Addr:   addr,
		State:  tv.StateDisconnected,
	}
	r.devices[id] = d
	return d, nil
}

// GetDevice gets a device by ID.
func (r *Registry) GetDevice(id string) (*Device, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.devices[id]
	return d, ok
}

// ListDevices lists all devices.
func (r *Registry) ListDevices() []*Device {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*Device, 0, len(r.devices))
	for _, d := range r.devices {
		res = append(res, d)
	}
	return res
}

// SetDriver sets the driver instance for a device.
func (d *Device) SetDriver(drv tv.Driver) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Drv = drv
	if drv != nil {
		d.State = tv.StateConnected
	} else {
		d.State = tv.StateDisconnected
	}
}

// GetDriver gets the driver instance.
func (d *Device) GetDriver() tv.Driver {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Drv
}

// SetState sets device state.
func (d *Device) SetState(s tv.State) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.State = s
}

// GetState gets device state.
func (d *Device) GetState() tv.State {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.State
}

// Call executes a function on the device's driver safely.
func (d *Device) Call(ctx context.Context, fn func(tv.Driver) error) error {
	drv := d.GetDriver()
	if drv == nil {
		return fmt.Errorf("driver not connected")
	}
	return fn(drv)
}
