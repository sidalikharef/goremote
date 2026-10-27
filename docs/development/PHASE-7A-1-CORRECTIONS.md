# Phase 7A.1 — Security Corrections

## Changes

### Safe Binding
- Default bind changed to `127.0.0.1:8080` (was `:8080`)
- LAN binding requires explicit `--addr` configuration
- Host validation improved with URL parsing

### Authentication
- Lightweight auth system added (`internal/auth`) with session cookies/tokens
- Credentials stored with 0600 perms, no hardcoded defaults
- State-changing APIs protected by RequireAuth
- Enrollment/auth endpoints need to be added for full flow (design documented)

### Origin/Host Validation
- Proper URL parsing with scheme/hostname/port validation
- Rejects malformed/cross-origin requests
- Host header not trusted as auth

### Security Headers
- CSP updated to avoid external scripts (removed unpkg)
- Local HTMX shim added for offline use

### Tests
- Added security tests for unauthorized access, cross-origin, invalid inputs
- All tests pass with race detector

## Security Implications of LAN Mode
- Binding to 0.0.0.0/LAN exposes unauthenticated ECP access risk
- Must use strong local auth, firewall rules, trusted network
- HTTPS strongly recommended if behind proxy

## Residual Risks
- Auth flow not fully wired in HTTP handlers yet
- CSRF protection needs SameSite cookies
- Rate limiting not implemented
