package postgres

import (
	"time"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

// Config holds PostgreSQL connection and pool settings.
type Config struct {
	Connection       Connection    `yaml:"connection"`
	Pool             *Pool         `yaml:"pool"`
	ConnectTimeout   time.Duration `yaml:"connect_timeout"`
	StatementTimeout time.Duration `yaml:"statement_timeout"`
}

// Connection holds PostgreSQL connection parameters.
type Connection struct {
	Hosts      string      `yaml:"hosts"`
	Username   string      `yaml:"username"`
	Password   string      `yaml:"password"`
	Database   string      `yaml:"database"`
	Additional string      `yaml:"additional"`
	TLS        *config.TLS `yaml:"tls"`
}

// Pool holds connection pool settings.
type Pool struct {
	MaxConnections        int32         `yaml:"max_connections"`
	MinConnections        int32         `yaml:"min_connections"`
	MaxConnectionLifetime time.Duration `yaml:"max_connection_lifetime"`
	MaxConnectionIdleTime time.Duration `yaml:"max_connection_idle_time"`
	HealthCheckPeriod     time.Duration `yaml:"health_check_period"`
}
