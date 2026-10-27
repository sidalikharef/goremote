# AGENTS.md — Project Conventions and Phase Rules

## Project Overview
GoRemote is an open-source, self-hosted universal TV remote. One Go binary, go:embed, standard-library-first, vanilla JS + HTMX. Local network only. No DB, no mandatory cloud, no Node runtime, no frontend frameworks.

## Code Conventions
- Go: follow standard Go style (gofmt). Use `go vet`. Keep packages small and focused.
- Interfaces in `tv/` define contracts; drivers implement them. Add optional capabilities via type assertions/interfaces to avoid breaking changes.
- All outbound HTTP to TVs must have timeouts (context + client). Validate all inputs.
- Errors: return wrapped errors with context; don't leak sensitive details to HTTP responses unnecessarily.
- Security: validate IPs/hosts/keys/app IDs; prevent SSRF/path traversal; serve embedded assets safely; default to localhost binding.
- Concurrency: be safe for concurrent handler/driver calls.
- No secrets in logs or repo.

## Architecture Rules
- Do not rewrite from scratch. Extend via interfaces and new packages.
- Standard library first. Avoid new deps unless absolutely necessary and justified.
- `go:embed` only for web assets.
- Keep public API contracts stable.
- Each driver self-contained under `drivers/<name>/`.

## Phase Execution Rules
- Phase 0: Audit/docs only. No code behavior changes. Create/update docs as specified.
- Phase 1+: Implement incrementally. Do not skip tests/mocks for hardware interactions.
- No git push, destructive commands, or commits unless explicitly requested.
- Do not claim unsupported protocols work.
- Stop at end of current phase and report blockers for review.

## Verification
- `go build ./...` must pass
- `go test ./...` must pass  
- `go vet ./...` clean
- `gofmt -l .` should show no files needing formatting

## File Locations
- Core interfaces: `tv/`
- Drivers: `drivers/`
- HTTP: `httpapi/` or handlers in main; keep organized
- Web: `web/` (embedded)
- Config: `config/`
- Docs: `docs/development/`
