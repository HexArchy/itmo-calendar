package config

type Config struct {
	App        *AppInfo    `path:"app"`
	Logger     *Logger     `path:"logger"`
	Shutdown   *Shutdown   `path:"shutdown"`
	HTTPServer *HTTPServer `path:"http_server"`
	Postgres   *Postgres   `path:"postgres"`
	RabbitMQ   *RabbitMQ   `path:"rabbitmq"`
	ITMO       *ITMO       `path:"itmo"`
	TLS        *TLS        `path:"tls"`
	Secrets    *Secrets    `path:"secret"`
}
