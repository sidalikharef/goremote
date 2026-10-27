package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	authpkg "github.com/sidalikharef/goremote/internal/auth"
)

// AuthServer handles auth endpoints
type AuthServer struct {
	auth *authpkg.Store
}

func NewAuth(a *authpkg.Store) *AuthServer {
	return &AuthServer{auth: a}
}

func (a *AuthServer) status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	hasUsers := a.auth.HasUsers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"has_users": hasUsers, "enrollment_available": !hasUsers})
}

func (a *AuthServer) enroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !authpkg.IsLoopback(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if a.auth.HasUsers() {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var req struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		BootstrapToken string `json:"bootstrap_token"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024*64)).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 || len(req.Password) > 1024 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := a.auth.Enroll(strings.TrimSpace(req.Username), req.Password, req.BootstrapToken); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *AuthServer) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024*64)).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !a.auth.Authenticate(strings.TrimSpace(req.Username), req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := a.auth.CreateSession()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "goremote_session",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   86400,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (a *AuthServer) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	token := authpkg.GetSessionToken(r)
	if token != "" {
		a.auth.RevokeSession(token)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "goremote_session",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusOK)
}
