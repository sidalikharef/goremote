FROM node:22 AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml* ./
RUN corepack enable && pnpm install --no-frozen-lockfile
COPY frontend/ ./
RUN pnpm build

FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN mkdir -p httpapi/web web && cp -r frontend/dist/* httpapi/web/ && cp -r frontend/dist/* web/ && rm -f httpapi/web/app.js httpapi/web/style.css web/app.js web/style.css
RUN go build ./...

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/goremote /app/goremote
EXPOSE 8080
ENTRYPOINT ["/app/goremote"]
