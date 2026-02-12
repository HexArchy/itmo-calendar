package http

import (
	"time"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

// Config holds HTTP server settings.
type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	TLS *config.TLS `yaml:"tls"`

	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	EnableHTTP2  bool          `yaml:"enable_http2"`
}
