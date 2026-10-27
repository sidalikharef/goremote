# Phase 7A.2 — Authentication and Security Completion

## Implemented
- Argon2id password hashing with crypto/rand salts
- First-run enrollment via CLI (`cmd/enroll`) with bootstrap token, loopback-only web enrollment path
- Login/logout/session status endpoints
- Session cookies HttpOnly, SameSite=Lax
- Auth middleware integrated, protects API endpoints
- Login/logout UI added
- Removed fake HTMX shim
- Proper origin/host validation maintained

## Usage
```bash
# Generate bootstrap and enroll (local)
go run ./cmd/enroll -username admin -password <pass>
# Start server (default localhost:8080)
go run .
```

## Security Notes
- Plain HTTP on LAN provides no confidentiality - use VPN/HTTPS proxy
- CSRF protection via SameSite cookies + origin checks
- No hardcoded credentials
- Sessions expire after 24h
