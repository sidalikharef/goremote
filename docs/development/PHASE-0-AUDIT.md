# Phase 0 — Repository Audit and Implementation Blueprint

**Date:** 2026-10-27  
**Scope:** Full repository audit (go.mod/main.go/packages/drivers/web/assets/tests/build/docs) against README as intended product doc.

## Repository Inventory
- Tracked files (HEAD): README.md only
- Working tree: README.md modified (link updated)
- Git history: single commit `9c9c1ae docs: add README`
- No Go files, no build scripts, no tests, no embedded assets, no package dirs present
- `.git/` contains standard git metadata; `opencode` marker present
- No go.mod, main.go, tv/, drivers/, web/ directories

Conclusion: the codebase as-is is *documentation-only*. The README describes a target architecture/implementation that does not exist in the repository yet.

## How It *Would* Work (from README)
- One Go binary, go:embed for web UI (web/index.html, web/style.css)
- HTTP server in main.go, flags/env config (--tv/--driver/--addr/--timeout), default driver roku
- tv.Driver interface: Name() string, Key(key string) error, Apps() ([]App, error), Launch(appID string) error
- drivers/roku/client.go implements Roku ECP (HTTP on :8060), XML for /query/apps and /query/device-info
- HTTP routes: /keypress/{key}, /query/apps, /launch/{id} implied; UI served from embedded assets
- HTMX + vanilla JS in UI

## Implementation Status Matrix (documented vs actual)
| Feature | Status | Evidence | Notes |
| --- | --- | --- | --- |
| Go module | missing | no go.mod | Cannot build |
| main.go / HTTP server / routing | missing | no files | Only described |
| tv.Driver interface | missing | no files | Only in README |
| drivers/roku (ECP client) | missing | no files | Only described |
| web/ with go:embed assets | missing | no files | Only described |
| Config flags/env | missing | no implementation | Only described |
| Keypress API/UI | missing | no implementation | Only described |
| Apps listing/launch | missing (partial in roadmap) | no implementation | Marked partial in README; missing here |
| Device discovery (SSDP) | missing | no implementation | Future |
| Tests | missing | no *_test.go files | No tests present |
| Build scripts (Makefile, justfile) | missing | none found | Only `go build` in README |
| Lint/typecheck scripts | missing | none found | `go vet`, `gofmt` mentioned in README dev section |
| AGENTS.md / CONTRIBUTING.md / LICENSE | missing | not present | Referenced in README/badges |
| docs/screenshot-* | missing | referenced as placeholders | Not present |

## Security & Reliability Findings (by severity)
Since there is no code to analyze, these are *pre-implementation risks* based on the intended design (treat as required controls):

**Critical**
- SSRF: if any future endpoint accepts arbitrary host/URL from client input (discovery/proxy/bridge), must strictly restrict to LAN targets, block loopback/metadata, enforce allowlists.
- Uncontrolled LAN requests: ECP has no auth; any process on LAN could be abused if HTTP server exposed beyond localhost or if endpoints are unauthenticated. Default bind to localhost or require explicit LAN bind with warning; consider mTLS/local-only by default.
- Path injection / traversal: serving embedded/static files must not allow `..` escaping embed roots; use `http.FileServer` with `embed.FS` carefully or serve specific files.
- CSRF: if browser-based UI posts state-changing requests from other origins, need same-site cookies/CSRF tokens or at least ensure endpoints are not GET for mutations and CORS locked down. Default to same-origin only.

**High**
- Input validation: key names, app IDs, IP addresses must be validated (allowlists for keys; numeric/uuid for app IDs; IP parsing). Avoid string concatenation into URLs/paths unsafely.
- Timeouts: all outbound HTTP to TVs must use per-request/timeouts (README states `--timeout` default 5s) - must enforce context timeouts, dial timeouts.
- Error handling: distinguish network errors, timeouts, TV errors; don't leak internal details to UI in unsafe ways.
- Concurrency: shared HTTP clients, driver instances accessed concurrently need synchronization; avoid global mutable state without locking.
- Pairing/persistence: future PIN pairing must store secrets securely (file perms, avoid world-readable). No plaintext cloud sync.

**Medium**
- DoS/resource limits: rate-limit keypresses? maybe not critical for LAN but reasonable.
- Logging: avoid logging full IPs/tokens in plaintext if added later.

## Architecture Evaluation (multi-TV, protocols, discovery, pairing)
Target design (README) is extensible via tv.Driver interface. Good separation. Needs:
- Registry/factory for drivers (by name), config store for saved TVs (multiple), selection/switching in UI
- Capability model (e.g., hasApps, hasVolume, hasPower) to hide UI elements
- Discovery abstraction (SSDP/mDNS/scan) with timeouts and dedup
- Pairing flow per protocol (LG webOS SSAP PIN, Android TV Remote v2 TLS+PIN)
- Persistent config (local file, no DB) with permissions
- HTTP API versioning/pluralization considerations
- Consistent error model

## Proposed Package Structure (incremental, non-breaking)
```
goremote/
├── main.go              # flag/env parsing, DI, server setup, routing
├── config/              # config loading/saving (flags+env+file), validation
├── tv/                  # core interfaces/types: Driver, App, Capabilities
├── drivers/
│   ├── registry.go      # driver registration/factory
│   └── roku/            # client.go, types.go, tests with mocks
├── httpapi/             # HTTP handlers, middleware (timeout, recovery, CORS, security headers)
├── web/                 # go:embed assets (index.html, style.css, app.js if added later)
└── internal/            # shared util (httpclient, xml, validate, netutil) - private
```
Constraints: standard library first, go:embed, no DB, no cloud, no new runtimes/frameworks.

## Driver Capability Design
Extend interface or use capabilities:
```go
type Driver interface {
  Name() string
  Key(key string) error
  Apps() ([]App, error)
  Launch(appID string) error
  // optional via type assertion or Capabilities()
}
type Capabilities struct {
  Apps bool; Volume bool; Power bool; TextInput bool; Discovery bool
}
func (d Driver) Caps() Capabilities { ... } // or separate interface
```
Keep backward compatible by adding optional interfaces (e.g., Discoverer, Pairer).

## Prioritized Phase Breakdown & Prerequisites
- Phase 0 (this): Audit/docs only. Prereq: none. Done when docs created.
- Phase 1: Bootstrap minimal working Roku implementation matching README quick start. Prereq: go.mod, interfaces, roku client, server, embed UI. Must not break API shape implied. Add mocks/tests.
- Phase 2: Universal/multi-TV, discovery, pairing, saved config.
- Phase 3: ESP32 IR bridge.

## Concrete Acceptance Criteria for Phase 1
- [ ] `go build ./...` succeeds with Go 1.22+
- [ ] `go test ./...` passes (including mocked Roku tests)
- [ ] `./goremote --tv 192.168.1.100` starts server on :8080, serves UI
- [ ] POST to /keypress/Home returns 200 and forwards to TV (mocked in tests)
- [ ] GET /query/apps returns parsed apps or appropriate error
- [ ] POST /launch/{id} works
- [ ] All outbound calls have timeouts; invalid inputs rejected with 4xx
- [ ] No secrets logged; binds configurable; embedded assets load
- [ ] README quick start remains accurate or updated only if necessary and compatible

## Verification Commands (Phase 0)
- `find . -type f -not -path '*/\.git/*' | sort`
- `git -C . log --oneline -5 && git -C . status`
- `go version` (if go present)

## Blockers / Notes
- No source code exists. Implementation must be created in subsequent phases per constraints. Phase 0 is doc-only.

## Additional Details

### Module/Go code
- **Status:** Missing. No `go.mod` found. Need to initialize module (e.g., `module github.com/sidalikharef/goremote` per README clone URL) in Phase 1.
- **Entry point:** `main.go` not present. Must handle flags/env/config, HTTP server setup, driver wiring.
- **Packages:** None exist yet.

### Interfaces/drivers
- **Status:** Missing. `tv/` package and `drivers/` not present. Target interface from README must be created in `tv/` with `Driver` interface and `App` struct.
- **Design:** Keep minimal for Phase 1 (Roku only). Add optional interfaces later for capabilities/discovery/pairing.

### Server/HTTP/UI/assets
- **Status:** Missing. No HTTP server, no routes, no `web/` assets, no go:embed usage.
- **Target:** Serve embedded `web/index.html`, `web/style.css` (HTMX). API endpoints as implied by README.

### Config
- **Status:** Missing. Flags/env as specified in README must be implemented.

### Tests/scripts/AGENTS/CONTRIBUTING/LICENSE/screenshots
- **Tests:** None. **Scripts:** None. **AGENTS.md:** Created in Phase 0. **CONTRIBUTING.md/LICENSE:** Referenced but not present. **Screenshots:** Referenced as placeholders, not present.

### Concrete Security Mitigations (to enforce in Phase 1+)
- **SSRF:** Never accept arbitrary URLs from clients; only talk to configured TV IP/host. Validate IP format, block private/loopback/metadata ranges only if exposing proxy (but we shouldn't). If discovery added, restrict to LAN.
- **LAN exposure/no auth:** Default bind `127.0.0.1` or `:8080` but document LAN usage; consider binding config. ECP has no auth - rely on trusted LAN and localhost-first.
- **Path traversal:** Use `embed.FS` and serve specific files or `http.FileServer(http.FS(webFS))` with care; reject paths containing `..`.
- **CSRF:** For same-origin UI with HTMX/forms, ensure state changes are POST/PUT/DELETE not GET; add SameSite cookies if sessions added later. No cookies planned now.
- **Input validation:** Allowlist for key names (Home, Back, Up, Down, Left, Right, Select, Play, Rev, Fwd, VolumeUp, VolumeDown, VolumeMute, Power, etc.). Validate app IDs as non-empty safe strings. Validate IPs with `net.ParseIP`.
- **Timeouts:** Use `*http.Client` with `Timeout` and per-request `context.Context` with timeout (respect `--timeout`).
- **Errors:** Wrap with `%w`, return generic errors to HTTP clients, log details server-side.
- **Concurrency:** Protect shared state with mutexes if any; drivers should be safe for concurrent calls.
- **Secrets:** Never log TV IPs in verbose error strings unnecessarily; no secrets in repo.
- **DoS:** Consider basic rate limiting later if exposed broadly (not required for Phase 1 localhost).

### Enhanced Acceptance Criteria for Phase 1
- [ ] `go mod init` creates valid module
- [ ] `tv` package defines `Driver` interface and `App` struct matching intended contract
- [ ] `drivers/roku` implements `tv.Driver`; uses `net/http` with timeouts; parses `/query/apps` XML safely
- [ ] `config` loads flags/env with defaults/validation
- [ ] `httpapi` serves UI from embed, handles keypress/apps/launch with validation
- [ ] `web/` contains `index.html` and `style.css` (HTMX) matching README UX
- [ ] Tests use `httptest` to mock Roku HTTP; table-driven where appropriate
- [ ] `go build ./...`, `go test ./...`, `go vet ./...` all pass
- [ ] `gofmt -l .` shows no unformatted files
- [ ] No new dependencies beyond stdlib
