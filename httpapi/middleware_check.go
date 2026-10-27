package httpapi

import (
	"net/http"
	"net/url"
	"strings"
)

// CheckOriginSafe uses proper URL parsing.
func CheckOriginSafe(allowedHosts map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			origin := r.Header.Get("Origin")
			if origin == "" {
				host := r.Host
				if allowedHosts != nil && allowedHosts[host] {
					next.ServeHTTP(w, r)
					return
				}
				w.WriteHeader(http.StatusForbidden)
				return
			}
			u, err := url.Parse(origin)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if u.Scheme != "http" && u.Scheme != "https" {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h := u.Hostname()
			p := u.Port()
			if p == "" {
				if u.Scheme == "http" {
					p = "80"
				} else {
					p = "443"
				}
			}
			hostPort := h + ":" + p
			if allowedHosts != nil && (allowedHosts[hostPort] || allowedHosts[u.Host]) {
				next.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusForbidden)
		})
	}
}

// CheckHostSafe uses proper validation.
func CheckHostSafe(allowed map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allowed == nil {
				next.ServeHTTP(w, r)
				return
			}
			h := r.Host
			if allowed[h] {
				next.ServeHTTP(w, r)
				return
			}
			// try to parse
			if strings.Contains(h, ":") {
				if allowed[h] {
					next.ServeHTTP(w, r)
					return
				}
			}
			w.WriteHeader(http.StatusBadRequest)
		})
	}
}
