# Contributing

Contributions are welcome! Here's how to get started.

## Development Setup

```bash
git clone https://github.com/HexArchy/itmo-calendar.git
cd itmo-calendar
make deps
make build
```

Requires: Go 1.24+, PostgreSQL, RabbitMQ (see `docker-compose.dev.yml`).

## Workflow

1. Fork the repo and create a branch from `master`
2. Make your changes
3. Run checks:
   ```bash
   make fmt
   make lint
   make test
   ```
4. Open a pull request against `master`

## Guidelines

- Follow existing code style and project structure
- Keep PRs focused — one feature or fix per PR
- Write tests for new functionality
- Update documentation if behavior changes

## Commit Messages

Use clear, concise commit messages:
- `Fix auth token refresh on expired session`
- `Add rate limiting to CalDAV endpoint`
- `Update RabbitMQ reconnect backoff logic`

## Code of Conduct

Be respectful and constructive. See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
