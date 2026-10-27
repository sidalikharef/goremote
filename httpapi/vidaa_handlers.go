package httpapi

import (
	"encoding/json"
	"net/http"

	authpkg "github.com/sidalikharef/goremote/internal/auth"
	"github.com/sidalikharef/goremote/internal/validate"
)

func (s *Server) handleVidaaManualRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if s.auth != nil {
		tok := authpkg.GetSessionToken(r)
		if tok == "" || !s.auth.ValidateSession(tok) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	// CSRF/origin handled by middleware
	var req struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Address string `json:"address"`
		Enable  bool   `json:"enable"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.ID == "" || req.Name == "" || req.Address == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if validate.IP(req.Address) != nil && req.Address == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// experimental disabled by default; require explicit enable
	if !req.Enable {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if _, err := s.reg.RegisterDevice(req.ID, req.Name, "vidaa", req.Address); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleVidaaPairSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if s.auth != nil {
		tok := authpkg.GetSessionToken(r)
		if tok == "" || !s.auth.ValidateSession(tok) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "pending"})
}

func (s *Server) handleVidaaPairStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if s.auth != nil {
		tok := authpkg.GetSessionToken(r)
		if tok == "" || !s.auth.ValidateSession(tok) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "pending"})
}

func (s *Server) handleVidaaPairCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if s.auth != nil {
		tok := authpkg.GetSessionToken(r)
		if tok == "" || !s.auth.ValidateSession(tok) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "cancelled"})
}

func (s *Server) handleVidaaUnpair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if s.auth != nil {
		tok := authpkg.GetSessionToken(r)
		if tok == "" || !s.auth.ValidateSession(tok) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "unpaired"})
}
