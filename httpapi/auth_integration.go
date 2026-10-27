package httpapi

import (
	"net/http"
)

// MountAuth mounts auth routes
func MountAuth(mux *http.ServeMux, as *AuthServer) {
	mux.HandleFunc("/api/auth/status", as.status)
	mux.HandleFunc("/api/auth/enroll", as.enroll)
	mux.HandleFunc("/api/auth/login", as.login)
	mux.HandleFunc("/api/auth/logout", as.logout)
}
