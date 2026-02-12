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
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("Starting application components")

			go func() {
				if err := srv.Start(); err != nil {
					l.Error("HTTP server failed", zap.Error(err))
				}
			}()

			go cronjob.New(
				cronUC,
				locker,
				rabbitCfg.Queues.CronProcessScheduleQueue,
				1*time.Minute,
				l.With(zap.String("component", "cron-scheduler")),
			).Start(ctx)

			go func() {
				if err := worker.Start(ctx); err != nil {
					l.Error("Send-schedule worker failed", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
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
