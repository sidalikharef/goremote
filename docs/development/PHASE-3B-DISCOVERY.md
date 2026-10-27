# Phase 3B — Real Device Discovery and Registration

## Implemented
- Roku SSDP discovery (ST: roku:ecp) with bounded timeouts/cancellation
- Discovery endpoints: POST /api/discover, GET /api/discover/list, POST /api/discover/save
- Registered devices listed via /api/devices
- Frontend discovery UI with explicit add action
- Devices persisted via config system
- No auto-save; approval required
- Invalid device ID does not fallback

## Security
- Auth required for mutations
- IP validation
- No uncontrolled subnet scans
- Untrusted discovery data treated carefully

## Tests
- Mock SSDP, registration, persistence, routing, invalid device rejection

## Hardware Tests Required
- Real Roku on LAN for SSDP
- Network conditions
