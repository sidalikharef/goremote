# GoRemote

GoRemote is an open-source, self-hosted universal TV remote. One Go binary, go:embed, standard-library-first, React + TypeScript frontend.

## Quick start

```bash
make build
./goremote
```

Open http://127.0.0.1:8080

## Frontend development

```bash
cd frontend
pnpm install
pnpm dev
```

Vite proxies API calls to `http://127.0.0.1:8080`.

## Build from scratch

```bash
make build-frontend
go build ./...
```

## Notes
- Local network only by default (binds to 127.0.0.1)
- Authentication required after first enrollment (loopback enrollment)
- Roku discovery supported. VIDAA is experimental and disabled by default.
