FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /go/bin/itmo-calendar ./cmd/itmo-calendar

FROM alpine:3.21

RUN addgroup -S app && adduser -S app -G app
RUN apk add --no-cache ca-certificates tzdata curl

COPY --from=builder /go/bin/itmo-calendar /usr/local/bin/itmo-calendar

RUN mkdir -p /etc/itmo-calendar && chown app:app /etc/itmo-calendar

USER app
WORKDIR /etc/itmo-calendar

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

ENTRYPOINT ["/usr/local/bin/itmo-calendar", "--config=/etc/itmo-calendar/config.yaml"]
