package app

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	httpserver "github.com/hexarchy/itmo-calendar/internal/handlers/http"
	sendschedule "github.com/hexarchy/itmo-calendar/internal/handlers/workers/send-schedule"
	"github.com/hexarchy/itmo-calendar/migrations"

	joblocker "github.com/hexarchy/itmo-calendar/internal/adapters/repositories/job-locker"
	preparesendschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/prepare-send-schedule"
	cronjob "github.com/hexarchy/itmo-calendar/pkg/cron-job"
	"github.com/hexarchy/itmo-calendar/pkg/logger"
	"github.com/hexarchy/itmo-calendar/pkg/postgres"
	"github.com/hexarchy/itmo-calendar/pkg/rabbitmq"
)

// Run registers lifecycle hooks that start HTTP server, cron scheduler and worker.
func Run(
	lc fx.Lifecycle,
	srv *httpserver.Server,
	worker *sendschedule.Worker,
	cronUC *preparesendschedule.UseCase,
	locker *joblocker.Repository,
	rabbitCfg *rabbitmq.Config,
	l *zap.Logger,
	shutdowner fx.Shutdowner,
) {
	var cronCtx context.Context
	var cronCancel context.CancelFunc
	var workerCtx context.Context
	var workerCancel context.CancelFunc

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("Starting application components")

			serverErrCh := make(chan error, 1)
			go func() {
				if err := srv.Start(); err != nil {
					l.Error("HTTP server failed to start", zap.Error(err))
					serverErrCh <- err
				}
			}()

			cronCtx, cronCancel = context.WithCancel(context.Background())
			cronRunner := cronjob.New(
				cronUC,
				locker,
				rabbitCfg.Queues.CronProcessScheduleQueue,
				1*time.Minute,
				l.With(zap.String("component", "cron-scheduler")),
			)
			go cronRunner.Start(cronCtx)

			workerCtx, workerCancel = context.WithCancel(context.Background())
			workerErrCh := make(chan error, 1)
			go func() {
				if err := worker.Start(workerCtx); err != nil {
					l.Error("Send-schedule worker failed to start", zap.Error(err))
					workerErrCh <- err
				}
			}()

			select {
			case err := <-serverErrCh:
				l.Error("HTTP server startup failed, shutting down app", zap.Error(err))
				_ = shutdowner.Shutdown()
				return err
			case err := <-workerErrCh:
				l.Error("Worker startup failed, shutting down app", zap.Error(err))
				_ = shutdowner.Shutdown()
				return err
			case <-time.After(100 * time.Millisecond):
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("Stopping application components")

			if cronCancel != nil {
				cronCancel()
			}
			if workerCancel != nil {
				workerCancel()
			}

			l.Info("Stopping HTTP server")
			if err := srv.Stop(ctx); err != nil {
				l.Error("HTTP server stop failed", zap.Error(err))
			}

			logger.Sync(l)
			return nil
		},
	})
}

// RunMigrations applies database migrations at startup.
func RunMigrations(cfg *postgres.Config, l *zap.Logger) error {
	l.Info("Migration – started")
	dbstring := postgres.BuildConnectionURI(cfg.Connection)

	if err := migrations.ApplyMigrations(context.Background(), l, dbstring); err != nil {
		return err
	}

	l.Info("Migration – completed successfully")
	return nil
}
