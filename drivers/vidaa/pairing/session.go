package pairing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type Session struct {
	ID        string
	DeviceID  string
	State     string
	ExpiresAt time.Time
	mu        sync.RWMutex
}

var (
	mu       sync.RWMutex
	sessions = map[string]*Session{}
)

func NewSession(deviceID string, ttl time.Duration) (*Session, error) {
	if deviceID == "" {
		return nil, errors.New("device required")
	}
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	s := &Session{
		ID:        hex.EncodeToString(b),
		DeviceID:  deviceID,
		State:     "pending",
		ExpiresAt: time.Now().Add(ttl),
	}
	mu.Lock()
	sessions[s.ID] = s
	mu.Unlock()
	return s, nil
}

func Get(id string) (*Session, bool) {
	mu.RLock()
	defer mu.RUnlock()
	s, ok := sessions[id]
	if !ok {
		return nil, false
	}
	if time.Now().After(s.ExpiresAt) {
		return nil, false
	}
	return s, true
}

func (s *Session) SetState(st string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.State = st
}

func (s *Session) Cancel() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.State = "cancelled"
}

func (s *Session) SubmitPIN(ctx context.Context, pin string) error {
	if len(pin) != 4 {
		return errors.New("invalid pin")
	}
	for _, c := range pin {
		if c < '0' || c > '9' {
			return errors.New("invalid pin")
		}
	}
	return nil
}
