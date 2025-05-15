package container

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/config"
)

type Container struct {
	Config *config.Config
	Logger *zap.Logger

	Infra    Infra
	Adapters Adapters
	Services Services
	UseCases UseCases
	Workers  Workers
}

func New(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*Container, error) {
	c := &Container{
		Config: cfg,
		Logger: logger,
	}

	err := c.initInfra(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "init infra")
	}

	err = c.initAdapters()
	if err != nil {
		return nil, errors.Wrap(err, "init repositories")
	}

	err = c.initServices()
	if err != nil {
		return nil, errors.Wrap(err, "init services")
	}

	err = c.initUseCases()
	if err != nil {
		return nil, errors.Wrap(err, "init use cases")
	}

	err = c.initWorkers()
	if err != nil {
		return nil, errors.Wrap(err, "init workers")
	}

	return c, nil
}
