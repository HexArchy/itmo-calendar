package config

import "time"

type Postgres struct {
	Connection PostgresConnection `path:"connection" desc:"database connection settings"`
	Pool       *PostgresPool      `path:"pool"`

	ConnectTimeout   time.Duration `default:"5s" desc:"database connection timeout"`
	StatementTimeout time.Duration `default:"30s" desc:"database statement timeout"`
}

type PostgresConnection struct {
	Hosts      string `default:"127.0.0.1:5432" desc:"database hosts in the format host:port,host:port"`
	Username   string `secret:"true" desc:"database username"`
	Password   string `secret:"true" desc:"database password"`
	Database   string `desc:"database name"`
	Additional string `default:"target_session_attrs=read-write" desc:"additional property"`
	TLS        *TLS   `path:"tls" desc:"TLS settings"`
}

type PostgresPool struct {
	MaxConnections        int32         `default:"3" desc:"database connection pool max connections"`
	MinConnections        int32         `default:"1" desc:"database connection pool min connections"`
	MaxConnectionLifetime time.Duration `default:"1h" desc:"database connection pool max connection lifetime"`
	MaxConnectionIdleTime time.Duration `default:"30m" desc:"database connection pool max connection idle time"`
	HealthCheckPeriod     time.Duration `default:"10s" desc:"database connection pool health check period"`
}
