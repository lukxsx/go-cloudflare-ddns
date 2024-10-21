FROM golang:1.22-alpine AS builder

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY . .

RUN go get

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o go-cloudflare-ddns
RUN chown appuser:appgroup /app/go-cloudflare-ddns

FROM scratch

COPY --from=builder /app/go-cloudflare-ddns /app/go-cloudflare-ddns
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER appuser

ENTRYPOINT ["/app/go-cloudflare-ddns"]
