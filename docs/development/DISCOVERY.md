# Device Discovery

## Overview
GoRemote supports automatic discovery via SSDP/UPnP and mDNS/DNS-SD where available, with manual IP entry as fallback.

## Requirements
- Bounded discovery duration and concurrency
- Retry, cancellation, sensible timeouts
- Deduplication by stable device identity (USN/UUID preferred over IP)
- Handle IP address changes and rediscovery
- Explicit user action required before saving new devices
- Network-scoped access only
- Do not perform uncontrolled full-subnet scanning
- Open port is not proof of compatibility
- Handle malformed responses safely

## Architecture
- `discovery/` types and manager for deduplication/state
- `discovery/ssdp/` SSDP M-SEARCH (best effort, UDP 239.255.255.250:1900)
- `discovery/mdns/` mDNS/DNS-SD (best effort, stdlib only)
- Integration with Phase 2 registry and config

## HTTP API
- `POST /api/discover` or `GET /api/discover` - start discovery (returns discovered candidates)
- `POST /api/discover/save` - save a discovered device (explicit user action required)
- `GET /api/devices` - list registered devices

## Security
- Validate all discovered data
- Only accept IP addresses, no arbitrary URLs
- Path traversal prevention
- No auto-pairing without confirmation
