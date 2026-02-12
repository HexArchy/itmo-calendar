package rabbitmq

import (
	"fmt"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

// Config holds RabbitMQ connection settings.
type Config struct {
	Host     string      `yaml:"host"`
	Port     int         `yaml:"port"`
	User     string      `yaml:"user"`
	Password string      `yaml:"password"`
	VHost    string      `yaml:"vhost"`
	TLS      *config.TLS `yaml:"tls"`
	Queues   *Queues     `yaml:"queues"`
}

// Queues holds queue names.
type Queues struct {
	CronProcessScheduleQueue string `yaml:"cron_process_schedule"`
	SendScheduleQueue        string `yaml:"send_schedule"`
}

// BuildDSN returns an AMQP connection string.
func (c *Config) BuildDSN() string {
	scheme := "amqp"
	if c.TLS != nil && c.TLS.Enabled {
		scheme = "amqps"
	}

	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s",
		scheme,
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.VHost,
	)
}
