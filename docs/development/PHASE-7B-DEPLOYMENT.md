# Phase 7B — Deployment, Distribution and Release Preparation

## Builds
- Single binary with embedded assets
- Cross-compilation supported via GOOS/GOARCH
- CI workflow added

## Deployment
- systemd unit for non-root operation
- Dockerfile and docker-compose (bind to localhost)
- Use HTTPS reverse proxy or VPN for remote access

## Security
- Default localhost binding
- Auth required
- No secrets in repo
