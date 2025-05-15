# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates build-base

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -extldflags '-static'" -o /go/bin/itmo-calendar ./cmd/itmo-calendar

# Runtime stage
FROM alpine:3.21

# Add non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create necessary directories with proper permissions
RUN mkdir -p /etc/itmo-calendar/certs /var/log/itmo-calendar \
    && chown -R appuser:appgroup /etc/itmo-calendar /var/log/itmo-calendar

# Copy the binary from builder
COPY --from=builder /go/bin/itmo-calendar /usr/local/bin/itmo-calendar
RUN chmod +x /usr/local/bin/itmo-calendar

# Copy configuration
COPY configs/itmo-calendar.docker.yaml /etc/itmo-calendar/config.yaml

# Set working directory
WORKDIR /etc/itmo-calendar

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8443

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/itmo-calendar", "--config=/etc/itmo-calendar/config.yaml"]
