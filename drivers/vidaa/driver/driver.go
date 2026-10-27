package driver

import (
	"context"
	"errors"
	"sync"

	"github.com/sidalikharef/goremote/drivers/vidaa/capabilities"
	"github.com/sidalikharef/goremote/drivers/vidaa/remote"
	"github.com/sidalikharef/goremote/drivers/vidaa/transport"
	"github.com/sidalikharef/goremote/tv"
)

type Driver struct {
	mu       sync.RWMutex
	tr       *transport.Client
	clientID string
	deviceID string
	auth     bool
	remote   *remote.Client
	vol      *capabilities.VolumeService
	svc      *capabilities.Service
}

func New(tr *transport.Client, clientID, deviceID string) *Driver {
	d := &Driver{
		tr:       tr,
		clientID: clientID,
		deviceID: deviceID,
		auth:     false,
	}
	d.remote = remote.New(tr, clientID, deviceID, false)
	d.vol = capabilities.NewVolume(tr, clientID, deviceID, false)
	d.svc = capabilities.New(tr, clientID, deviceID, false)
	return d
}

func (d *Driver) Name() string { return "vidaa" }

func (d *Driver) SetAuthenticated(auth bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.auth = auth
	d.remote = remote.New(d.tr, d.clientID, d.deviceID, auth)
	d.vol = capabilities.NewVolume(d.tr, d.clientID, d.deviceID, auth)
	d.svc = capabilities.New(d.tr, d.clientID, d.deviceID, auth)
}

func (d *Driver) IsAuthenticated() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.auth
}

func (d *Driver) Key(ctx context.Context, key string) error {
	d.mu.RLock()
	auth := d.auth
	d.mu.RUnlock()
	if !auth {
		return errors.New("not paired")
	}
	return d.remote.SendKey(ctx, key)
}

func (d *Driver) Apps(ctx context.Context) ([]tv.App, error) {
	d.mu.RLock()
	auth := d.auth
	d.mu.RUnlock()
	if !auth {
		return nil, errors.New("not paired")
	}
	return nil, errors.New("not implemented in adapter")
}

func (d *Driver) Launch(ctx context.Context, appID string) error {
	d.mu.RLock()
	auth := d.auth
	d.mu.RUnlock()
	if !auth {
		return errors.New("not paired")
	}
	return d.svc.LaunchApp(ctx, appID)
}

func (d *Driver) Capabilities() tv.Capabilities {
	return tv.Capabilities{
		Navigation: true,
		Volume:     true,
		Power:      true,
		Media:      true,
		Apps:       true,
		TextInput:  false,
		Inputs:     true,
		Pairing:    true,
	}
}
