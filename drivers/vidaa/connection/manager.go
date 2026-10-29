package connection

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"
	"time"

	vidaaCreds "github.com/sidalikharef/goremote/drivers/vidaa/credentials"
	vidaaDrv "github.com/sidalikharef/goremote/drivers/vidaa/driver"
	"github.com/sidalikharef/goremote/drivers/vidaa/transport"
)

type State string

const (
	StateDisconnected State = "disconnected"
	StateConnecting   State = "connecting"
	StateConnected    State = "connected"
	StatePairing      State = "pairing_required"
	StateError        State = "error"
)

type Manager struct {
	mu    sync.RWMutex
	state State
	drv   *vidaaDrv.Driver
	tr    *transport.Client
	store *vidaaCreds.Store
}

func New(store *vidaaCreds.Store) *Manager {
	return &Manager{store: store, state: StateDisconnected}
}

func (m *Manager) State() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

func (m *Manager) Connect(ctx context.Context, addr string, clientID, deviceID string) error {
	m.mu.Lock()
	m.state = StateConnecting
	m.mu.Unlock()
	if addr == "" {
		m.mu.Lock()
		m.state = StateError
		m.mu.Unlock()
		return errors.New("address required")
	}
	broker := "tls://" + addr + ":36669"
	tlsCfg := &tls.Config{MinVersion: tls.VersionTLS12}
	tr, err := transport.New(transport.Config{Broker: broker, ClientID: clientID, TLS: tlsCfg, Timeout: 10 * time.Second})
	if err != nil {
		m.mu.Lock()
		m.state = StateError
		m.mu.Unlock()
		return err
	}
	if err := tr.Connect(ctx); err != nil {
		m.mu.Lock()
		m.state = StateError
		m.mu.Unlock()
		return err
	}
	m.tr = tr
	m.drv = vidaaDrv.New(tr, clientID, deviceID)
	if m.store != nil {
		if c, err := m.store.Load(deviceID); err == nil && c != nil && c.Token != "" {
			m.drv.SetAuthenticated(true)
			m.mu.Lock()
			m.state = StateConnected
			m.mu.Unlock()
			return nil
		}
	}
	m.mu.Lock()
	m.state = StatePairing
	m.mu.Unlock()
	return nil
}

func (m *Manager) Driver() *vidaaDrv.Driver     { return m.drv }
func (m *Manager) Transport() *transport.Client { return m.tr }

func (m *Manager) Disconnect() {
	if m.tr != nil {
		m.tr.Disconnect(250)
	}
	m.mu.Lock()
	m.state = StateDisconnected
	m.mu.Unlock()
}
