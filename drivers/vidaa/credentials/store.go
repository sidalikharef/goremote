package credentials

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

type Creds struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Token        string `json:"token,omitempty"`
	DeviceID     string `json:"device_id"`
}

type Store struct {
	path string
	mu   sync.RWMutex
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) Save(c Creds) error {
	if c.DeviceID == "" {
		return errors.New("device required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *Store) Load(deviceID string) (*Creds, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var c Creds
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	if c.DeviceID != deviceID {
		return nil, errors.New("device mismatch")
	}
	return &c, nil
}

func (s *Store) Delete() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return os.Remove(s.path)
}
