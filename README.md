# ITMO Calendar

A CalDAV service for ITMO University scheduling system.

## Quick Start with Docker Compose

```bash
# Clone the repository
git clone https://github.com/hexarchy/itmo-calendar.git
cd itmo-calendar

# Create environment file
cat > .env.prod << EOF
# ITMO Calendar Production Environment
POSTGRES_PASSWORD=your_secure_password
RABBITMQ_PASSWORD=your_secure_password
JWT_SECRET=your_secure_jwt_secret
EOF

# Generate certificates if needed
./generate-certs.sh

# Start the services
docker-compose --env-file .env.prod -f docker-compose.prod.yml up -d
```

## Docker Image

Docker image is available on GitHub Container Registry:

```bash
docker pull ghcr.io/hexarchy/itmo-calendar:latest
```
