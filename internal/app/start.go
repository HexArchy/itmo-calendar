package app

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	cronjob "github.com/hexarchy/itmo-calendar/pkg/cron-job"
	"github.com/hexarchy/itmo-calendar/pkg/shutdown"
)

// Start launches all application components and handles graceful shutdown.
func (a *App) Start(ctx context.Context) error {
	a.Logger.Info("Starting ITMO calendar",
		zap.String("environment", a.Cfg.App.Environment),
		zap.String("version", a.Cfg.App.Version))

	err := a.startMigrations(ctx)
	if err != nil {
		return errors.Wrap(err, "start migrations")
	}

	// Initialize runners for concurrent component startup
	runners := make(map[string]func(context.Context) error)

	// HTTP server runner
	runners["http"] = func(_ context.Context) error {
		a.Logger.Info("Starting HTTP server")
		err := a.HTTPServer.Start()
		if err != nil {
			return errors.Wrap(err, "start HTTP server")
		}
		return nil
	}

	runners["cron-scheduler"] = func(ctx context.Context) error {
		a.Logger.Info("Starting cron scheduler")
		runner := cronjob.New(a.Container.UseCases.PrepareSendSchedule,
			a.Container.Adapters.JobLocker,
			a.Cfg.RabbitMQ.Queues.CronProcessScheduleQueue,
			1*time.Minute, // TODO: make it configurable.
			a.Logger.With(zap.String("component", "cron-scheduler")),
		)
		runner.Start(ctx)
		a.Logger.Info("Cron scheduler started")
		return nil
	}

	runners["send-schedule"] = func(ctx context.Context) error {
		a.Logger.Info("Starting workers")
		err := a.Container.Workers.RabbitMQ.SendSchedule.Start(ctx)
		if err != nil {
			return errors.Wrap(err, "start workers")
		}
		a.Logger.Info("Workers started")
		return nil
	}

	// Start all registered runners.
	errCh := make(chan error, len(runners))
	var wg sync.WaitGroup
	for name, runFn := range runners {
		// Create local copy for goroutine.
		name := name
		runFn := runFn
		a.Logger.Info("Starting component", zap.String("component", name))
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := runFn(ctx)
			if err != nil {
				a.Logger.Error("Component failed",
					zap.String("component", name),
					zap.Error(err))
				errCh <- errors.Wrapf(err, "%s failed", name)
				shutdown.Shutdown()
			}
		}()
	}

	// Set up detection for shutdown.
	shutdown.AddCallback(&shutdown.Callback{
		Name: "postgres connection",
		FnCtx: func(_ context.Context) error {
			if a.Container.Infra.Postgres != nil {
				a.Container.Infra.Postgres.Close()
				a.Logger.Info("PostgreSQL connection closed")
			}
			return nil
		},
	})

	shutdown.AddCallback(&shutdown.Callback{
		Name: "HTTP server",
		FnCtx: func(ctx context.Context) error {
			err := a.HTTPServer.Stop(ctx)
			if err != nil {
				return errors.Wrap(err, "stop HTTP server")
			}
			return nil
		},
	})

	shutdown.AddCallback(&shutdown.Callback{
		Name: "RabbitMQ connection",
		FnCtx: func(ctx context.Context) error {
			err := a.Container.Infra.RabbitMQ.Close()
			if err != nil {
				return errors.Wrap(err, "stop RabbitMQ connection")
			}
			return nil
		},
	})

	shutdownStarted := make(chan struct{})
	go func() {
		for {
			if shutdown.IsShuttingDown() {
				close(shutdownStarted)
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(100 * time.Millisecond):
			}
		}
	}()

	select {
	case <-ctx.Done():
		a.Logger.Info("Application context canceled, initiating shutdown")
		shutdown.Shutdown()
	case err := <-errCh:
		return err
	case <-shutdownStarted:
		a.Logger.Info("Shutdown signal received")
	}

	config := &shutdown.Config{
		Delay:           a.Cfg.Shutdown.Delay,
		WaitTimeout:     a.Cfg.Shutdown.Timeout,
		CallbackTimeout: a.Cfg.Shutdown.CallbackTimeout,
	}

	a.Logger.Info("Waiting for graceful shutdown to complete...")
	err = shutdown.Wait(config)
	if err != nil {
		a.Logger.Error("Failed to gracefully shutdown application", zap.Error(err))
		return err
	}

	a.Logger.Info("Application gracefully stopped")
	return nil
}
