# Phase 5A — Hisense VIDAA Integration

## Stage A — Protocol Investigation (Completed)

### Research findings (public/community)
- VIDAA OS is Hisense's smart TV platform. Common local control approaches vary by generation/firmware.
- Some VIDAA versions expose HTTP/REST APIs (often on non-standard ports) or use MQTT/other protocols; documentation is fragmented.
- Others rely on mobile app APIs; not all are openly documented.
- Do **not** assume Android TV ADB/Remote v2 or webOS SSAP.
- Do **not** assume standard ports without evidence.

### Discovery
- Possible mDNS/SSDP in some firmwares, but not guaranteed. Must not rely on open-port probing alone.
- Network-scoped discovery only.

### Authentication/Pairing
- Some implementations require PIN/pairing or tokens; varies. Must store credentials securely (0600 perms), per Phase 2 config model.
- Treat as unknown until verified.

### Feasibility Matrix
| Aspect | Status | Evidence | Notes |
| --- | --- | --- | --- |
| Confirmed local protocol spec | Unknown/Not established | Public docs incomplete; vendor docs not present | Cannot proceed to implementation safely |
| Confirmed default ports/endpoints | Unknown | No verified references in repo | Do not guess |
| Confirmed auth/pairing | Unknown | No verified references | Must not invent |
| Confirmed key commands | Unknown | No verified references | Must not claim support |
| Model/firmware compatibility | Unknown | No target model specified | Need device info |
| Discovery method | Unclear | Possibly mDNS/SSDP in some versions | Not guaranteed |

### Supported vs Unsupported (pre-implementation)
- **Supported:** None claimed.
- **Unsupported:** All operations until verified.

## Stage B — Implementation
**Not performed.** Protocol not sufficiently established.

## Required Device Information / Diagnostic Evidence
To proceed safely to Stage B, provide:
1. **Target device:** Model number (e.g., Hisense 55U7G, specific VIDAA version) and firmware version
2. **Network evidence:** Packet capture (pcap) of official mobile app controlling TV, or HAR/logs of app traffic
3. **Protocol details:** Confirmed endpoints (host:port/path), methods, headers, auth flow, message formats
4. **Pairing:** Step-by-step pairing procedure and how credentials are obtained/stored
5. **Commands:** Verified list of commands that work (power, nav, volume, apps) with examples
6. **Discovery:** Confirmed mDNS/SSDP service names/types if used
7. **Constraints:** Known version differences

## Next Steps
STOP after Stage A as instructed. No VIDAA driver implemented. No new TV platforms added.
