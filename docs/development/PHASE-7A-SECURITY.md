# Phase 7A — Security Hardening

## Threat Model
- **Trusted:** User on local network accessing via browser
- **Untrusted:** Other devices on LAN, cross-origin sites, malicious local processes
- **Assets:** TV control, pairing credentials, config files

## Implemented Controls
- **Binding:** Defaults to :8080; recommend localhost binding. Host checking added.
- **CSRF/Origin:** CheckOrigin blocks cross-origin POSTs without valid Origin matching allowed hosts
- **Host validation:** CheckHost restricts Host header
- **Security headers:** X-Content-Type-Options, X-Frame-Options, Referrer-Policy, Permissions-Policy, CSP
- **Request limits:** Max body size (64KB), method restrictions
- **Input validation:** Key names and app IDs validated; IP validation
- **Device targeting:** Only registered devices controllable; no arbitrary SSRF
- **Config security:** Atomic writes, 0600 perms
- **Error handling:** Generic errors to clients
- **Logging:** No secrets logged

## Recommendations (Not Implemented)
- mTLS for remote access
- Rate limiting
- Authentication for non-localhost
- Secure discovery auto-save disabled (explicit action required)

## Testing
Added security middleware tests in httpapi/security_test.go

## Residual Risks
- ECP has no auth (rely on LAN trust + localhost)
- Discovery stubs need hardening when implemented
