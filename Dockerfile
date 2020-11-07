FROM golang:1.15 as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -a -v -o /app/qbittorrent-notify

FROM alpine:3.10
RUN apk --no-cache add \
  ca-certificates
COPY --from=builder /app/qbittorrent-notify /usr/local/bin

RUN chmod +x /usr/local/bin/qbittorrent-notify

ENTRYPOINT ["/usr/local/bin/qbittorrent-notify"]
