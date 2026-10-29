package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/argon2"
)

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"` // hex argon2id
	Salt         string `json:"salt"`
	Version      int    `json:"version"`
}

type Store struct {
	mu        sync.RWMutex
	users     map[string]User
	path      string
	sessions  map[string]session
	bootstrap string
}

type session struct {
	expires time.Time
}

func New(path string) *Store {
	s := &Store{
		users:     make(map[string]User),
		path:      path,
		sessions:  make(map[string]session),
		bootstrap: path + ".bootstrap",
	}
	s.load()
	return s
}

func (s *Store) load() {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		// fail closed behavior - don't load corrupt data
		return
	}
	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return
	}
	for _, u := range users {
		if u.Version == 1 { // only support new version
			s.users[u.Username] = u
		}
	}
}

func (s *Store) save() error {
	s.mu.RLock()
	users := make([]User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	s.mu.RUnlock()
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	if err := os.Chmod(filepath.Dir(s.path), 0700); err != nil {
	}
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func hashPassword(password, salt string) (string, error) {
	saltB, err := hex.DecodeString(salt)
	if err != nil {
		return "", err
	}
	h := argon2.IDKey([]byte(password), saltB, 3, 64*1024, 4, 32)
	return hex.EncodeToString(h), nil
}

func (s *Store) GenerateBootstrapToken() (string, error) {
	if s.HasUsers() {
		return "", errors.New("users exist")
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)
	if err := os.WriteFile(s.bootstrap, []byte(token), 0600); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Store) CheckBootstrap(token string) bool {
	if token == "" {
		return false
	}
	data, err := os.ReadFile(s.bootstrap)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(token), data) == 1
}

func (s *Store) ClearBootstrap() error {
	return os.Remove(s.bootstrap)
}

func (s *Store) HasUsers() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users) > 0
}

func (s *Store) Enroll(username, password, bootstrapToken string) error {
	if s.HasUsers() {
		return errors.New("enrollment disabled")
	}
	if !s.CheckBootstrap(bootstrapToken) {
		return errors.New("invalid bootstrap")
	}
	if len(username) < 1 || len(password) < 8 {
		return errors.New("invalid credentials")
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	saltStr := hex.EncodeToString(salt)
	h, err := hashPassword(password, saltStr)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.users[username] = User{Username: username, PasswordHash: h, Salt: saltStr, Version: 1}
	s.mu.Unlock()
	if err := s.save(); err != nil {
		return err
	}
	return s.ClearBootstrap()
}

func (s *Store) Authenticate(username, password string) bool {
	s.mu.RLock()
	u, ok := s.users[username]
	s.mu.RUnlock()
	if !ok {
		return false
	}
	h, err := hashPassword(password, u.Salt)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(h), []byte(u.PasswordHash)) == 1
}

func (s *Store) CreateSession() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)
	s.mu.Lock()
	s.sessions[token] = session{expires: time.Now().Add(24 * time.Hour)}
	s.mu.Unlock()
	return token, nil
}

func (s *Store) ValidateSession(token string) bool {
	if token == "" {
		return false
	}
	s.mu.RLock()
	sess, ok := s.sessions[token]
	s.mu.RUnlock()
	if !ok {
		return false
	}
	if time.Now().After(sess.expires) {
		s.mu.Lock()
		delete(s.sessions, token)
		s.mu.Unlock()
		return false
	}
	return true
}

func (s *Store) RevokeSession(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func GetSessionToken(r *http.Request) string {
	c, err := r.Cookie("goremote_session")
	if err == nil && c != nil {
		return c.Value
	}
	return r.Header.Get("X-GoRemote-Token")
}

func IsLoopback(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsLoopback() {
		return true
	}
	return false
}
