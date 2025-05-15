package config

import "time"

type HTTPServer struct {
	Host string `path:"host" default:"localhost"`
	Port int    `path:"port" default:"8080"`

	// TLS settings.
	TLS *TLS `path:"tls" desc:"TLS settings"`

	// Advanced server settings.
	ReadTimeout  time.Duration `path:"read_timeout" default:"5s"`
	WriteTimeout time.Duration `path:"write_timeout" default:"5s"`
	IdleTimeout  time.Duration `path:"idle_timeout" default:"60s"`
	EnableHTTP2  bool          `path:"enable_http2" default:"true"`
}
