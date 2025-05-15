package container

import (
	"context"

	"github.com/hexarchy/itmo-calendar/pkg/rabbitmq"

	"github.com/pkg/errors"
)

func (c *Container) initRabbitMQ(ctx context.Context) (*rabbitmq.Client, error) {
	tls, err := c.Config.RabbitMQ.TLS.BuildTLSConfig(c.Config.RabbitMQ.Host)
	if err != nil {
		return nil, errors.Wrap(err, "init rabbitmq tls config")
	}

	rabbitMQ, err := rabbitmq.New(ctx, c.Config.RabbitMQ.BuildDSN(), tls, c.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "init rabbitmq client")
	}

	return rabbitMQ, nil
}
