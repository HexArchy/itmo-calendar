# ITMO Calendar

[![CI](https://github.com/HexArchy/itmo-calendar/actions/workflows/ci.yml/badge.svg)](https://github.com/HexArchy/itmo-calendar/actions/workflows/ci.yml)
[![Deploy](https://github.com/HexArchy/itmo-calendar/actions/workflows/deploy.yml/badge.svg)](https://github.com/HexArchy/itmo-calendar/actions/workflows/deploy.yml)

Syncs your ITMO University schedule to any calendar app. Authenticates via ITMO SSO, fetches schedule from `my.itmo.ru`, converts to iCal and serves as a subscribable calendar feed.

> **Production:** [itmo-calendar.duckdns.org](https://itmo-calendar.duckdns.org)

---

## Add to your calendar

### Apple Calendar (iPhone / Mac) — CalDAV

1. **Settings** > **Calendar** > **Accounts** > **Add Account** > **Other** > **Add CalDAV Account**
2. Enter:
   - **Server:** `itmo-calendar.duckdns.org`
   - **Username:** your ISU number (e.g. `123456`)
   - **Password:** your ITMO password
3. Done. Calendar "ITMO Schedule" will appear and refresh automatically.

### Google Calendar / Android

- **Android:** use [DAVx5](https://www.davx5.com/) → Add CalDAV account → `https://itmo-calendar.duckdns.org/caldav/` with ISU + password
- **Desktop:** Thunderbird, GNOME Calendar — add CalDAV account with the same URL

### Any CalDAV client

Server: `https://itmo-calendar.duckdns.org`
CalDAV path: `/caldav/` (auto-discovery via `/.well-known/caldav`)
Auth: HTTP Basic (`ISU:password`)

### iCal subscription (legacy)

For clients that only support iCal subscriptions (not CalDAV):

URL: `https://itmo-calendar.duckdns.org/cal`
Auth: HTTP Basic (`ISU:password`)

---

## API

Base URL: `https://itmo-calendar.duckdns.org`

| Method | Endpoint | Description |
|--------|----------|-------------|
| `PROPFIND` | `/caldav/{isu}/` | CalDAV principal (Basic Auth) |
| `PROPFIND` | `/caldav/{isu}/calendars/schedule/` | CalDAV calendar collection |
| `GET` | `/.well-known/caldav` | CalDAV auto-discovery redirect |
| `GET` | `/cal` | iCal feed (Basic Auth) |
| `GET` | `/api/v1/health` | Health check |
| `POST` | `/api/v1/subscribe` | Subscribe by ISU + password |
| `GET` | `/api/v1/{isu}/ical` | Get iCal for user |
| `GET` | `/api/v1/{isu}/schedule` | Get schedule as JSON |
| `GET` | `/docs` | Interactive API docs (Scalar) |
| `GET` | `/openapi.yaml` | OpenAPI spec |

---

## Self-hosting

```bash
git clone https://github.com/HexArchy/itmo-calendar.git
cd itmo-calendar

# Configure secrets
cat > .env <<EOF
POSTGRES_PASSWORD=changeme
RABBITMQ_PASSWORD=changeme
JWT_SECRET=changeme
EOF

# Run
docker compose up -d
```

Image: `ghcr.io/hexarchy/itmo-calendar:latest`

### Architecture

```
Client ──> Chi router ──> ogen API / CalDAV handler
                │
         ┌──────┼──────────┐
         │      │          │
     PostgreSQL RabbitMQ  ITMO API
         │      │
     migrations workers (cron + send-schedule)
```

Go 1.24 · Chi v5 · ogen · pgx · fx · amqp091-go

---

## Development

```bash
make build          # Build binary
make run            # Build + run locally
make test           # Tests with race detector
make lint           # staticcheck + go vet
make fmt            # gofmt + goimports
make swagger-gen    # Regenerate ogen code
```

## License

MIT
