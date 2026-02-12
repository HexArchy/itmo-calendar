# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| latest  | Yes       |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public issue
2. Email: [open a private security advisory](https://github.com/HexArchy/itmo-calendar/security/advisories/new)
3. Include steps to reproduce, impact assessment, and suggested fix if possible

You can expect an initial response within 48 hours.

## Scope

- Authentication flow (ITMO SSO / OAuth PKCE)
- Credential storage and token handling
- HTTP endpoints and input validation
- Docker image and deployment configuration

## Out of Scope

- ITMO University infrastructure (`my.itmo.ru`, `id.itmo.ru`)
- Third-party dependencies (report upstream)
