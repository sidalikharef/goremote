# Phase 7A.3 — Final Authentication and Integration Repair

## Authentication
- Auth routes bypassed in middleware correctly
- Enrollment requires loopback + bootstrap token, one-time only
- Argon2id with bounded params, password length limits
- Sessions HttpOnly/SameSite, revocation
- CLI enrollment reads password securely (no plaintext in args)

## Integration
- Frontend sends device ID explicitly
- API rejects unknown device IDs (no fallback)
- Auth UI integrated, logout/expired handled
- Discovery/manual add marked unavailable
- Removed fake HTMX shim

## Tests (httptest + mocks)
- Fresh install status
- Unauthorized command blocked
- Enrollment flow (first success, second rejected)
- Login/logout, wrong password
- Authorized mock Roku command
- Invalid device ID not falling back
- All pass with -race

## Verification
All tests pass. go build/vet/fmt/mod verify clean.
