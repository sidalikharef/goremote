# Roadmap

## Phase 0 — Repository Audit and Implementation Blueprint (current)
- [x] Read repo, inspect structure (none existing)
- [x] Document how intended system works from README
- [x] Create implementation-status matrix
- [x] Identify security/reliability issues
- [x] Propose architecture/package structure/driver capabilities
- [x] Prioritize phases and acceptance criteria
- [x] Create PHASE-0-AUDIT.md, ARCHITECTURE.md, ROADMAP.md
- [x] Update/create AGENTS.md with conventions/phase rules

## Phase 1 — Roku driver (minimal working)
Prerequisites: go.mod, tv interface, roku client, HTTP server, embedded UI, config
- [ ] Module init (go.mod), basic structure
- [ ] `tv/Driver` interface + types (`App`)
- [ ] `drivers/roku` client: Key, Apps (parse XML), Launch; timeouts, validation
- [ ] `httpapi` handlers: /, /keypress/{key}, /query/apps, /launch/{id} (or adjust to match README)
- [ ] `web/` embed index.html/style.css (HTMX) per README design
- [ ] Config flags/env (--tv, --driver, --addr, --timeout)
- [ ] Tests with httptest/mocks for Roku HTTP
- [ ] Manual verification against README quick start shape
Acceptance: see PHASE-0-AUDIT.md

## Phase 2 — Universal remote (multi-TV, discovery, pairing)
- [ ] Driver registry/factory, multi-TV config store (local file)
- [ ] Capabilities model, UI switching
- [ ] Discovery (SSDP/mDNS) with timeouts
- [ ] Samsung Tizen (WebSocket :8002) - only if implemented with care
- [ ] LG webOS (SSAP :3001, PIN pairing)
- [ ] Android TV/Google TV (Remote v2 TLS, PIN pairing)
- [ ] Secure pairing storage

## Phase 3 — ESP32 IR bridge
- [ ] `irbridge` driver
- [ ] Firmware design notes
- [ ] IR learn/persist

## Notes
- Do not break public API contracts
- Standard library first, go:embed, no DB/cloud/frameworks
- No unsupported protocols claimed as working

## Phase 3 — Automatic device discovery (current)
- [x] Discovery subsystem with SSDP/UPnP and mDNS/DNS-SD support (where applicable)
- [x] Safe device identification, deduplication by stable ID
- [x] IP change/rediscovery handling
- [x] Manual IP fallback
- [x] Bounded duration, concurrency, retry/cancel/timeouts
- [x] Explicit user action before saving
- [x] Network-scoped validation, no full-subnet scan
- [x] HTTP endpoints for discovery/registry
- [x] Mock discovery tests and duplicate-device tests
- [x] Malformed response handling
- [x] Documentation updated
