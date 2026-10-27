package httpapi

import (
	"net/http"
	"strings"
)

// SecurityHeaders adds common security headers.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'")
		next.ServeHTTP(w, r)
	})
}

// MaxBody limits request body size.
func MaxBody(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AllowedMethods restricts methods.
func AllowedMethods(methods ...string) func(http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, m := range methods {
		allowed[m] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(methods) > 0 && !allowed[r.Method] {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CheckOrigin validates Origin for state-changing requests.
func CheckOrigin(allowedHosts map[string]bool) func(http.Handler) http.Handler {
	return CheckOriginSafe(allowedHosts)
}

// CheckHost validates Host header.
func CheckHost(allowed map[string]bool) func(http.Handler) http.Handler {
	return CheckHostSafe(allowed)
}

// RequireAuth requires authentication for state-changing/protected operations.
func RequireAuth(store interface {
	ValidateSession(string) bool
}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/" {
				next.ServeHTTP(w, r)
				return
			}
			if strings.HasPrefix(r.URL.Path, "/api/auth/") {
				next.ServeHTTP(w, r)
				return
			}
			if r.Method == http.MethodGet {
				if strings.HasSuffix(r.URL.Path, ".css") || strings.HasSuffix(r.URL.Path, ".js") || strings.HasSuffix(r.URL.Path, ".svg") || strings.HasSuffix(r.URL.Path, ".png") || strings.HasSuffix(r.URL.Path, ".json") || strings.HasSuffix(r.URL.Path, ".ico") {
					next.ServeHTTP(w, r)
					return
				}
			}
			token := getToken(r)
			if token == "" || store == nil || !store.ValidateSession(token) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getToken(r *http.Request) string {
	if c, err := r.Cookie("goremote_session"); err == nil && c != nil {
		return c.Value
	}
	return r.Header.Get("X-GoRemote-Token")
}
