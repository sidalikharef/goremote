FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o goremote .

FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/goremote /goremote
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/goremote"]
