# Architecture

## Goals
- One Go executable, `go:embed` for web assets
- Standard library first
- Vanilla JS + HTMX UI
- Interface-driven drivers per TV brand/protocol
- Local-network only by design; no mandatory cloud/DB
- Incremental, non-breaking changes

## Current State
Documentation-only repo (see PHASE-0-AUDIT.md). Target structure from README not present.

## Target Components
- `main.go`: CLI flags/env, DI, HTTP server, route wiring, graceful shutdown
- `tv/`: Core types and `Driver` interface
- `drivers/`: Registry + per-driver packages (e.g., `roku/`, future `tizen/`, `webos/`, `androidtv/`, `irbridge/`)
- `httpapi/`: Handlers, middleware (logging/recovery/timeout/security headers/CORS), input validation
- `config/`: Configuration model (flags/env/file), validation, defaults
- `web/`: Embedded UI assets
- `internal/`: Shared utilities (http client with timeouts, net/validation, errors)

## Driver Interface (as specified)
```go
package tv

type App struct {
    ID   string
    Name string
}

type Driver interface {
    Name() string
    Key(key string) error
    Apps() ([]App, error)
    Launch(appID string) error
}
```
Evolution: add optional interfaces (`Discoverer`, `Pairer`, `Power`, `Volume`) and `Capabilities()` to avoid breaking existing implementations.

## Request Flow
1. Browser loads embedded UI (GET /)
2. HTMX/JS posts to API (e.g., POST /keypress/Home)
3. `httpapi` validates input, calls `tv.Driver.Key(...)`
4. Driver uses HTTP client with timeouts to talk to TV (e.g., Roku ECP :8060)
5. Errors mapped to appropriate HTTP status codes; UI shows feedback

## Security Boundaries
- Default bind to localhost; explicit LAN bind allowed with caution
- Only allow known key names via allowlist
- Validate IPs/hosts; never proxy arbitrary user-supplied URLs to arbitrary destinations (SSRF)
- Timeouts on all outbound connections
- No logging of secrets
- Serve embedded files safely (no path traversal)
- State-changing methods should not be GET; consider CSRF protections if ever exposing cross-origin
- ECP has no auth - treat as trusted LAN only

## Concurrency
- HTTP server concurrent by default; driver methods may be called concurrently
- Shared clients must be safe; avoid global mutable state

## Configuration
Flags override env; file config for multi-TV/pairing in later phases. No DB required.

## Testing Strategy
- Unit tests with interface mocks/fakes for drivers
- HTTP handler tests with `httptest`
- Table-driven tests
- Hardware-dependent tests behind build tags or skipped in CI; document manual verification
