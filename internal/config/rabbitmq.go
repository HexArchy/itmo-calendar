package config

import (
	"fmt"
)

type RabbitMQ struct {
	Host     string  `path:"host" default:"localhost" desc:"RabbitMQ host"`
	Port     int     `path:"port" default:"5672" desc:"RabbitMQ port"`
	User     string  `path:"user" default:"guest" desc:"RabbitMQ username"`
	Password string  `path:"password" secret:"true" desc:"RabbitMQ password"`
	VHost    string  `path:"vhost" default:"/" desc:"RabbitMQ virtual host"`
	TLS      *TLS    `path:"tls" desc:"TLS settings"`
	Queues   *Queues `path:"queues" desc:"RabbitMQ queues"`
}

type Queues struct {
	CronProcessScheduleQueue string `path:"cron_process_schedule" default:"cron_process_schedule" desc:"RabbitMQ cron process schedule queue"`
	SendScheduleQueue        string `path:"send_schedule" default:"send_schedule" desc:"RabbitMQ send schedule queue"`
}

func (r *RabbitMQ) BuildDSN() string {
	scheme := "amqp"
	if r.TLS != nil && r.TLS.Enabled {
		scheme = "amqps"
	}

	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s",
		scheme,
		r.User,
		r.Password,
		r.Host,
		r.Port,
		r.VHost,
	)
}
