#!/bin/bash
# generate-certs.sh - Generate TLS certificates for ITMO Calendar

set -e

# Configuration
CERTS_DIR="./certs"
APP_CERTS_DIR="$CERTS_DIR"
PG_CERTS_DIR="$CERTS_DIR/postgres"
RABBITMQ_CERTS_DIR="$CERTS_DIR/rabbitmq"
CA_KEY="$CERTS_DIR/ca.key"
CA_CERT="$CERTS_DIR/ca.crt"
DAYS_VALID=365
DOMAIN=${DOMAIN:-81.31.244.102}
COUNTRY=${COUNTRY:-RU}
STATE=${STATE:-Saint-Petersburg}
LOCALITY=${LOCALITY:-Saint-Petersburg}
ORGANIZATION=${ORGANIZATION:-HexArch}
ORGANIZATIONAL_UNIT=${ORGANIZATIONAL_UNIT:-IT Department}
EMAIL=${EMAIL:-nikitabelekov@gmail.com}

# Create directories
mkdir -p "$APP_CERTS_DIR" "$PG_CERTS_DIR" "$RABBITMQ_CERTS_DIR"

# Generate Certificate Authority key and certificate
echo "Generating CA certificate..."
openssl genrsa -out "$CA_KEY" 4096
openssl req -x509 -new -nodes -key "$CA_KEY" -sha256 -days $DAYS_VALID -out "$CA_CERT" \
  -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/OU=$ORGANIZATIONAL_UNIT/CN=ITMO Calendar CA/emailAddress=$EMAIL"

# Function to generate certificates for a service
generate_cert() {
  local SERVICE=$1
  local CN=$2
  local CERT_DIR=$3
  local SERVER_KEY="$CERT_DIR/server.key"
  local SERVER_CSR="$CERT_DIR/server.csr"
  local SERVER_CERT="$CERT_DIR/server.crt"
  
  echo "Generating certificate for $SERVICE..."
  
  # Create OpenSSL config for SAN support
  cat > "$CERT_DIR/openssl.cnf" << EOF
[req]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
C = $COUNTRY
ST = $STATE
L = $LOCALITY
O = $ORGANIZATION
OU = $ORGANIZATIONAL_UNIT
CN = $CN
emailAddress = $EMAIL

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = $CN
DNS.2 = localhost
IP.1 = 127.0.0.1
IP.2 = 81.31.244.102
EOF

  # Generate private key
  openssl genrsa -out "$SERVER_KEY" 2048
  chmod 600 "$SERVER_KEY"
  
  # Generate CSR with SAN
  openssl req -new -key "$SERVER_KEY" -out "$SERVER_CSR" \
    -config "$CERT_DIR/openssl.cnf"
  
  # Create extension config file
  cat > "$CERT_DIR/v3.ext" << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = $CN
DNS.2 = localhost
IP.1 = 127.0.0.1
EOF

  # Generate server certificate signed by CA
  openssl x509 -req -in "$SERVER_CSR" -CA "$CA_CERT" -CAkey "$CA_KEY" \
    -CAcreateserial -out "$SERVER_CERT" -days $DAYS_VALID \
    -extfile "$CERT_DIR/v3.ext"
  
  # Copy CA certificate to the service directory
  if [ "$CA_CERT" != "$CERT_DIR/ca.crt" ]; then
    cp "$CA_CERT" "$CERT_DIR/ca.crt"
  fi
  
  # Clean up
  rm "$SERVER_CSR" "$CERT_DIR/openssl.cnf" "$CERT_DIR/v3.ext"
  
  echo "$SERVICE certificate generated successfully"
}

# Generate certificates for each service
generate_cert "ITMO Calendar App" "$DOMAIN" "$APP_CERTS_DIR"
generate_cert "PostgreSQL" "postgres" "$PG_CERTS_DIR" 
generate_cert "RabbitMQ" "rabbitmq" "$RABBITMQ_CERTS_DIR"

# Generate client certificate for mutual TLS if needed
if [ "$GENERATE_CLIENT_CERT" = "true" ]; then
  echo "Generating client certificate..."
  CLIENT_KEY="$APP_CERTS_DIR/client.key"
  CLIENT_CSR="$APP_CERTS_DIR/client.csr"
  CLIENT_CERT="$APP_CERTS_DIR/client.crt"
  CLIENT_P12="$APP_CERTS_DIR/client.p12"
  
  openssl genrsa -out "$CLIENT_KEY" 2048
  chmod 600 "$CLIENT_KEY"
  
  openssl req -new -key "$CLIENT_KEY" -out "$CLIENT_CSR" \
    -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/OU=Client/CN=client@$DOMAIN/emailAddress=$EMAIL"
  
  openssl x509 -req -in "$CLIENT_CSR" -CA "$CA_CERT" -CAkey "$CA_KEY" \
    -CAcreateserial -out "$CLIENT_CERT" -days $DAYS_VALID
  
  # Create P12 file for browser import
  openssl pkcs12 -export -out "$CLIENT_P12" -inkey "$CLIENT_KEY" -in "$CLIENT_CERT" \
    -certfile "$CA_CERT" -passout pass:changeit
  
  echo "Client certificate generated successfully"
  rm "$CLIENT_CSR"
fi

# Set permissions
find "$CERTS_DIR" -type f -name "*.key" -exec chmod 600 {} \;
find "$CERTS_DIR" -type f -name "*.crt" -exec chmod 644 {} \;


echo "All certificates generated successfully"
