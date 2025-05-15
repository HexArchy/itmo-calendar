package container

import (
	"context"

	"github.com/hexarchy/itmo-calendar/pkg/rabbitmq"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Infra struct {
	Postgres *pgxpool.Pool
	RabbitMQ *rabbitmq.Client
}

func (c *Container) initInfra(ctx context.Context) error {
	var err error

	c.Infra.Postgres, err = c.initPostgresEngine(ctx)
	if err != nil {
		return errors.Wrap(err, "init postgres engine")
	}

	c.Infra.RabbitMQ, err = c.initRabbitMQ(ctx)
	if err != nil {
		return errors.Wrap(err, "init rabbitmq client")
	}

	return nil
}
