# GoRemote

GoRemote is an open-source, self-hosted universal TV remote. One Go binary with `go:embed`, standard-library-first backend, and a React + TypeScript + Vite frontend. Local network only. No DB, no mandatory cloud, no Node runtime for serving.

## Features

- **Local-first**: Binds to localhost by default; designed for your local network.
- **Session-based auth**: Cookie/session authentication with secure password hashing (Argon2id).
- **First-time enrollment**: CLI-based enrollment from loopback with a one-time bootstrap token.
- **Login-first UX**: Anonymous users are redirected to the login screen.
- **Device discovery**: SSDP discovery for supported TVs (e.g. Roku).
- **Manual device registration**: Add TVs by IP address.
- **React UI**: Modern remote control UI, device management, discovery, and settings.
- **Experimental VIDAA support**: Hisense VIDAA support exists but is unverified and disabled by default.

## Supported and experimental platforms

| Platform | Status | Notes |
| --- | --- | --- |
| Roku | Supported | SSDP discovery and basic remote commands. |
| Hisense VIDAA | Experimental (disabled by default) | Research/in-progress; manual registration only, not guaranteed to work on all hardware. |

## Requirements

- Go 1.22+
- Node.js + pnpm (for building the frontend)
- Linux/macOS/Windows

## Installation from GitHub

```bash
git clone https://github.com/sidalikharef/goremote.git
cd goremote
```

## Building the React frontend and Go binary

```bash
# Build frontend
cd frontend
pnpm install --frozen-lockfile
pnpm build
cd ..

# Copy build to embedded assets (already done in repo, but rebuild if needed)
# The Go server embeds httpapi/web/ via go:embed

# Build Go binary
go build ./...
# or
make build
```

## Running GoRemote

By default, GoRemote listens on `127.0.0.1:8080`.

```bash
go run .
# or
./goremote
```

## First-time admin enrollment

GoRemote uses CLI enrollment to create the first admin account. The CLI derives its auth store from the config path (default `~/.config/goremote/config.json` and `~/.config/goremote/config.json.auth.json`).

1. Create config directory with secure permissions:
```bash
mkdir -p ~/.config/goremote
chmod 700 ~/.config/goremote
```

2. Enroll the first admin user:
```bash
go run ./cmd/enroll -username <your-username>
```
Replace `<your-username>` with your desired username. You can also specify `-config /path/to/config.json`.

What happens:
- If no users exist, the CLI generates a **bootstrap token** (one-time, used only for initial enrollment).
- You will be prompted for a **password** (input is hidden). Minimum length is 8 characters.
- On success, the CLI prints `enrolled successfully`.

**Important notes:**
- The bootstrap token is for initial enrollment only. It is not your login password.
- The password is what you use to sign in via the web UI.
- Session cookies are used for browser authentication after login.
- Do **not** use `cmd/maketok` for production enrollment (it writes to `/tmp/test.auth.json` for development/testing only).
- If users already exist, enrollment is disabled for security.

## Admin login

1. Start the server: `go run .` (or `./goremote`)
2. Open [http://127.0.0.1:8080/login](http://127.0.0.1:8080/login)
3. Sign in with the username and password you enrolled.

On first run with no users, the login page will show a "First-time setup" form if enrollment is available from loopback; otherwise use the CLI enrollment flow above.

## Device discovery

1. Go to the **Discovery** page in the UI (or navigate to `/discovery`).
2. Click **Scan for TVs**. GoRemote will scan the local network for supported devices (Roku via SSDP).
3. Discovered devices will appear as cards. Click **Add TV** to save them.
4. If nothing is found and your TV is off or on a different network, you'll see "No TVs found".

Discovery requires an authenticated session.

## Manual device registration

1. Go to **My TVs** (`/devices`).
2. Under "Add TV manually", enter a name, IP address, and driver (Roku by default).
3. Click **Add device** to register it.
4. Select the device from the header dropdown to control it from the Remote page.

For experimental VIDAA, use the "Experimental VIDAA" section on the same page (manual registration only; unverified).

## Docker usage

Docker support is not currently included in this repository. Run directly with Go/pnpm as described above.

## Development commands

```bash
# Frontend
cd frontend
pnpm install
pnpm dev        # Vite dev server with proxy
pnpm build      # Production build
pnpm typecheck
pnpm test

# Backend
go build ./...
go vet ./...
go test ./... -count=1
go test -race ./...
gofmt -l .
```

## Troubleshooting

- **Can't access the UI**: Make sure you're connecting to `http://127.0.0.1:8080` (default). The server binds to localhost only by default.
- **Login page doesn't appear / always redirected**: Ensure cookies are enabled and you're using the same origin. Check that no reverse proxy is interfering with cookies/headers.
- **Discovery finds nothing**: TV must be powered on and on the same local network. Some networks isolate devices (AP isolation). Try again.
- **"users already exist" on enroll**: An admin already exists. You cannot re-enroll from CLI; manage accounts via the existing auth store or remove the auth file only if you understand the consequences.
- **Origin/Host errors**: GoRemote enforces Host/Origin checks for security. Use the default localhost address or configure allowed hosts appropriately.

## Security notes

- Default bind is localhost only.
- All protected API routes require an authenticated session (cookie-based).
- Host and Origin headers are checked.
- Passwords are hashed with Argon2id.
- Session cookies are HttpOnly with SameSite=Lax.
- Auth files and bootstrap tokens are stored with restrictive permissions (0600/0700).
- No secrets are logged.

## Screenshots

Screenshots will be added as the UI stabilizes. See the UI in the repo under `frontend/src/` and built assets in `httpapi/web/`.

## Roadmap

- Broader TV brand support
- Improved discovery
- Better device capabilities detection
- Hardened VIDAA support (experimental)

## License

See [LICENSE](LICENSE).
