# Next Steps (Audit & Plan Only - No Implementation)

## Recommendation: Phase 7 (Security, Deployment, Release Readiness) OR strengthen Phase 1-4 hardening

**Preferred safe next phase:** **Phase 7 — Security, deployment and release readiness**

### Reasoning
- Core architecture (Phases 0-4) is largely implemented. Roku works via mocks.
- Discovery infrastructure exists.
- Android TV (5D) is properly BLOCKED pending verified evidence (PCAP/proto/certs). Forcing it would risk security violations (TLS bypass, guessing).
- VIDA A (5A) also BLOCKED pending evidence.
- Phase 7 adds value without requiring new TV hardware/protocol reverse engineering: security hardening, deployment docs, release process, CI, binary signing considerations, safer defaults, etc.
- Can be done with existing codebase.

**Alternative:** If hardware becomes available for Roku/others, focus on hardware verification. But no hardware specified.

## Immediate Prerequisites (Before Implementation)
- [ ] Review security controls against current implementation
- [ ] Add/verify .gitignore for secrets/config
- [ ] Review CORS, headers, binding defaults
- [ ] Add deployment docs (systemd, Docker, reverse proxy)
- [ ] Add release checklist and versioning
- [ ] Consider CI workflow (lint/test/build)

## Remaining Tasks for Phase 7
- Security audit pass (SSRF, traversal, CSRF, secrets)
- Deployment documentation
- Configuration examples
- Release readiness checklist
- Maybe add `go:generate` or build tags
- Ensure no sensitive logging
- Review file permissions

## Acceptance Criteria
- Security checklist documented and addressed
- Deployment guide added
- Release process documented
- All existing tests still pass
- No new deps unless justified
- Defaults remain localhost-safe

## Expected Code Areas
- `httpapi/` (middleware/headers), `main.go` (defaults), `config/` (perms), docs

## Testing Strategy
- Existing tests must pass
- Security-focused review

## Blockers
- None critical; can proceed if approved

## Dependencies on Hardware
- None for Phase 7

## Android TV Blocker Analysis
- **Established:** Protocol requires TLS + certificate-based pairing + protobuf; cannot bypass pairing/TLS verification.
- **Uncertain:** Exact ports/message schemas/firmware differences.
- **Existing specs:** Public Android TV Remote Protocol v2 exists in community, but implementation details vary; still need device-specific validation and secure cert handling.
- **Isolation:** Android TV can remain BLOCKED without blocking Phase 7. No dependency.

## Phase 7A Complete
Security hardening added (headers, origin/host checks, body limits). Tests added. Baseline maintained.
