#!/bin/bash
# deploy.sh - Deploy ITMO Calendar to production

set -e

# Configuration
ENV_FILE=".env.prod"
DOCKER_COMPOSE="docker-compose.prod.yml"

# Check if environment file exists
if [ ! -f "$ENV_FILE" ]; then
  echo "Creating environment file..."
  cat > "$ENV_FILE" << EOF
# ITMO Calendar Production Environment
POSTGRES_PASSWORD=$(openssl rand -base64 32)
RABBITMQ_PASSWORD=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 64)
EOF
  echo "Environment file created. Please review $ENV_FILE before continuing."
  exit 1
fi

# Check if certificates exist
if [ ! -d "./certs" ] || [ ! -f "./certs/server.crt" ]; then
  echo "TLS certificates not found. Generating certificates..."
  ./generate-certs.sh
fi

# Build and start containers
echo "Deploying ITMO Calendar services..."
docker-compose --env-file "$ENV_FILE" -f "$DOCKER_COMPOSE" up -d --build 


echo "Deployment complete! ITMO Calendar is running at https://81.31.244.102/api/v1"
