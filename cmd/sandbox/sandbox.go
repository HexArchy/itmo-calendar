package main

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/app/container"
	"github.com/hexarchy/itmo-calendar/internal/config"
)

type Sandbox struct {
	Cfg       *config.Config
	Logger    *zap.Logger
	Container *container.Container
	Tools     *Tools
}

func NewSandbox(ctx context.Context, cfg *config.Config) (*Sandbox, error) {
	s := Sandbox{
		Cfg: cfg,
	}

	var err error
	s.Logger = initLogger()

	s.Container, err = container.New(ctx, cfg, s.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "new container")
	}

	s.Tools = initTools()

	return &s, nil
}

func initLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	logger.Info("Init Logger – success")
	return logger
}

type Tools struct {
	// Add any tools needed for sandbox testing.
}

func initTools() *Tools {
	return &Tools{}
}
