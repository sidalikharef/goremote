# Phase 4.5 — Reproducible Baseline and Verification

## Environment
- OS: Fedora Linux 44 (Workstation Edition)
- Go: go1.26.10-X:nodwarf5 linux/amd64 (found at /usr/bin/go)
- go.mod: github.com/sidalikharef/goremote, go 1.22
- go mod verify: all modules verified

## Git Safety
- Branch: main (up to date with origin/main)
- Modified: README.md
- Untracked: AGENTS.md, config/, discovery/, docs/, drivers/, go.mod, httpapi/, internal/, main.go, tv/, web/
- No secrets/certs/keys detected in tracked/untracked source. Config written with 0600 perms. No .gitignore changes made. Files remain untracked (not staged/committed).

## Build/Test Results
```bash
go mod verify  # PASS
go build ./... # PASS (after fixes)
go vet ./...   # PASS (after fixes)
go test ./...  # PASS
go test -race ./... # PASS
```

Fixes made to establish baseline:
- `discovery/manager.go`: removed unused `context` import
- `httpapi/handlers.go`: fixed `registry.Device.Call` usage to match `func(tv.Driver) error`, removed duplicate `Routes()`, added discovery endpoints
- `discovery/ssdp/ssdp.go`: removed unused `remote` variable
- `discovery/types.go`: added IP validation in `Validate()` (prevents accepting non-IP addresses)
- `drivers/registry/registry_test.go`: removed unused variable
- Formatted with `gofmt`

## Feature Verification
| Feature | Status | Evidence |
| --- | --- | --- |
| Roku ECP navigation/keys | IMPLEMENTED_UNTESTED | `drivers/roku/client.go` - Key with validation/timeouts/context |
| Roku app launching/query | IMPLEMENTED_UNTESTED | Apps() XML parse, Launch() with validation |
| TV registry + capabilities | IMPLEMENTED_AND_TESTED | `drivers/registry/` + tests |
| Atomic config persistence | IMPLEMENTED_AND_TESTED | `config/` with atomic save, 0600 perms, tests |
| SSDP/mDNS discovery | STUB/PARTIAL | SSDP does UDP M-SEARCH (best effort), mDNS returns empty; manager/dedup tested |
| Discovery HTTP handlers | STUB | `/api/discover`, `/api/discover/save` return stubs (explicit save required) |
| Device selection | IMPLEMENTED_UNTESTED | Registry-based selection in handlers |
| Browser UI + embedded assets | IMPLEMENTED_UNTESTED | `web/` and `httpapi/web/` embedded via go:embed |
| PWA manifest/icons | IMPLEMENTED_UNTESTED | manifest.json, icon.svg present |
| Service worker | MISSING | No service worker found (intentionally not added per Phase 4 constraints) |
| Pairing contracts/cred storage | IMPLEMENTED_UNTESTED (interfaces) | `tv.Pairer` interface, config stores pairing with 0600 |

## Missing Tests
- HTTP API integration tests with fake Roku server
- SSDP/mDNS with mock UDP
- Discovery save flow
- Main wiring tests

## Security Notes
- No TLS bypass, no pairing bypass
- Input validation, timeouts, path traversal checks
- Config perms 0600
- Defaults appear localhost-friendly

## Hardware Verification Required
Roku ECP, discovery on real LAN, pairing flows - all unverified without hardware.

## Readiness for Phase 7A
Baseline compiles and tests pass. Ready for Phase 7A (Security/Deployment/Release) upon approval.
