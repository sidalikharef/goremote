# 📺 GoRemote

> A self-hosted, cross-platform TV remote that runs in your browser. One Go binary, no cloud, no app store, no JavaScript framework.

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

GoRemote is a local-network remote control for your TV. It talks directly to your TV's native control protocol over Wi-Fi — no extra hardware, no vendor app, no internet round-trip. Open `http://<your-pc>:8080` on your phone and you have a remote.

**Currently supported:** Roku
**In development:** Samsung Tizen, LG webOS, Android TV, and an ESP32 IR bridge for non-smart TVs.

---

## ✨ Why this exists

- **No app required.** It's a web page. Bookmark it on your phone's home screen.
- **No cloud.** Everything stays on your LAN. Your TV never talks to a third party.
- **One binary.** The web UI is compiled into the executable with `go:embed`. Copy one file to a Raspberry Pi and you're done.
- **No JS framework.** Vanilla HTML + [HTMX](https://htmx.org). The UI loads instantly and works on a 10-year-old phone.
- **Hackable.** Each TV brand is a small, isolated Go package behind a common interface. Adding a TV is a `client.go` file, not a rewrite.

---

## 📸 Screenshots

| Desktop | Mobile |
|---|---|
| _placeholder — add `docs/screenshot-desktop.png`_ | _placeholder — add `docs/screenshot-mobile.png`_ |

---

## 🚀 Quick Start

### Requirements

- Go 1.22 or newer
- A Roku on the same Wi-Fi network as the machine running GoRemote
- Roku **Settings → System → Advanced system settings → Control by mobile apps** set to **Enabled** (or **Network access: Permissive**)

### Install from source

```bash
git clone https://github.com/sidalikharef/goremote.git
cd goremote
go build -o goremote .
./goremote --tv 192.168.1.100
```

Then open **http://localhost:8080** in your browser.

### Find your TV's IP

On the TV: **Settings → Network → About**. Or scan your LAN:

```bash
# macOS / Linux
arp -a | grep -i roku

# or with nmap
nmap -p 8060 --open 192.168.1.0/24
```

Port `8060` is the Roku ECP port — if it's open, that's your TV.

---

## ⚙️ Configuration

GoRemote takes flags and reads environment variables.

| Flag | Env var | Default | Description |
|---|---|---|---|
| `--tv` | `GOREMOTE_TV_IP` | _(required)_ | IP address of the TV |
| `--driver` | `GOREMOTE_DRIVER` | `roku` | TV driver: `roku` (more soon) |
| `--addr` | `GOREMOTE_ADDR` | `:8080` | Address the web server listens on |
| `--timeout` | `GOREMOTE_TIMEOUT` | `5s` | Per-request timeout when talking to the TV |

Example:

```bash
GOREMOTE_TV_IP=192.168.1.100 GOREMOTE_ADDR=:9000 ./goremote
```

---

## 🗺️ Roadmap

- [x] **Phase 1 — Roku driver.** ECP client, web UI, single-binary build.
  - [x] Keypress API (`/keypress/{key}`)
  - [x] Web remote UI (D-pad, volume, playback, home/back)
  - [x] Embedded static assets
  - [ ] App launcher (`/query/apps` + `/launch/{id}`)
  - [ ] Text entry (`/keypress/Lit_` for on-screen keyboard)
  - [ ] Device auto-discovery via SSDP
- [ ] **Phase 2 — Universal remote.**
  - [ ] Samsung Tizen driver (WebSocket, port 8002)
  - [ ] LG webOS driver (SSAP, port 3001, PIN pairing)
  - [ ] Android TV / Google TV driver (Remote v2 over TLS, PIN pairing)
  - [ ] Driver auto-detect + saved config file
  - [ ] Multiple TVs in one UI, switchable
- [ ] **Phase 3 — ESP32 IR bridge.**
  - [ ] Firmware: web server + IR blaster on ESP32
  - [ ] Firmware: IR code learning from an existing remote
  - [ ] `irbridge` Go driver that talks to the ESP32
  - [ ] Persist learned codes to flash (SPIFFS)

Vote on what gets built next by opening an issue with the `roadmap` label.

---

## 🧱 Architecture

```
goremote/
├── main.go              # flag parsing, HTTP server, route wiring
├── tv/
│   └── driver.go        # Driver interface every TV brand implements
├── drivers/
│   └── roku/
│       └── client.go    # Roku ECP implementation
└── web/                 # embedded with //go:embed
    ├── index.html
    └── style.css
```

The entire design is built around one small interface:

```go
package tv

type Driver interface {
    Name() string
    Key(key string) error
    Apps() ([]App, error)
    Launch(appID string) error
}
```

The HTTP layer never knows what brand of TV it's talking to. Each driver package is self-contained and testable in isolation. **To add a TV, implement `Driver` and register it.** That's it.

---

## 🔌 How the Roku driver works

Roku exposes the **External Control Protocol (ECP)**: a plain HTTP server on port `8060` on your LAN.

| Purpose | Request |
|---|---|
| Press a key | `POST /keypress/Home` |
| List installed apps | `GET /query/apps` → XML |
| Launch an app | `POST /launch/12` |
| Query device info | `GET /query/device-info` → XML |

Key names include `Home`, `Back`, `Up`, `Down`, `Left`, `Right`, `Select`, `Play`, `Rev`, `Fwd`, `VolumeUp`, `VolumeDown`, `VolumeMute`, and `Power`.

No authentication is required, which is why this driver is the starting point for the project.

---

## 🛠️ Development

```bash
# run with live reload (needs air: go install github.com/air-verse/air@latest)
air

# build
go build -o goremote .

# test
go test ./...

# vet + format
go vet ./... && gofmt -l .
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the driver-authoring guide.

---

## ❓ FAQ

**Does this work when I'm not home?**
No, by design. GoRemote only talks to your TV on your local network. If you want remote access, put it behind a VPN (WireGuard/Tailscale) — don't port-forward it.

**Does it work on older Roku models?**
ECP has shipped on every Roku for over a decade. It's very likely fine. If it isn't, open an issue with your model number.

**Why Go and not Rust/Node?**
Single static binary, trivial cross-compilation to a Raspberry Pi, and a standard library that covers everything this project needs.

**Why no React?**
The entire UI is ~15 buttons. HTMX does everything needed in one `<script>` tag. Fewer dependencies means it still builds in five years.

**Can you add support for <brand>?**
Yes — that's the whole roadmap. Open an issue with the protocol details if you know them, or send a PR implementing `tv.Driver`.

---

## 🙏 Acknowledgments

- [shreve/go-roku](https://github.com/shreve/go-roku) — Go Roku ECP reference
- [Roku ECP documentation](https://developer.roku.com/docs/developer-program/debugging/external-control-api.md)
- [HTMX](https://htmx.org)
- [IRremoteESP8266](https://github.com/crankyoldgit/IRremoteESP8266) — for the upcoming ESP32 phase

---

## 📄 License

[MIT](LICENSE) © 2026 Sid Ali Kharef

**Not affiliated with Roku, Samsung, LG, or Google.** All product names are trademarks of their respective owners and are used here only to describe interoperability.

---

## 👤 About the Author

I'm **Sid Ali Kharef**, a developer building open-source tools in Go and C++. GoRemote is part of an ongoing effort to write software I actually want to use, and to learn network protocols by implementing them from scratch instead of wrapping a library.

If you're hiring, collaborating, or just want to talk shop about local-network protocols, embedded C++, or Go — my inbox is open.

<p align="left">
  <a href="https://www.linkedin.com/in/sid-ali-kharef-b5195419b"><img alt="LinkedIn" src="https://img.shields.io/badge/LinkedIn-Sid%20Ali%20Kharef-0A66C2?logo=linkedin&logoColor=white"></a>
  <a href="https://sidalikharef.cc"><img alt="Website" src="https://img.shields.io/badge/Website-sidalikharef.cc-000000?logo=firefox&logoColor=white"></a>
  <a href="https://github.com/sidalikharef"><img alt="GitHub" src="https://img.shields.io/badge/GitHub-sidalikharef-181717?logo=github&logoColor=white"></a>
  <a href="https://stackoverflow.com/users/12624124/sidali"><img alt="Stack Overflow" src="https://img.shields.io/badge/Stack%20Overflow-12624124-F58025?logo=stackoverflow&logoColor=white"></a>
  <a href="mailto:sidalikharef@gmail.com"><img alt="Email" src="https://img.shields.io/badge/Email-sidalikharef@gmail.com-EA4335?logo=gmail&logoColor=white"></a>
</p>

> 💼 **Open to work** — remote or on-site, backend / systems / embedded roles. [Resume (PDF)](https://sidalikharef.cc/files/Sid-Ali-Kharef-CV-Updated.pdf)