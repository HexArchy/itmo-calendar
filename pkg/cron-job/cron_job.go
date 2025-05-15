package cronjob

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Job interface {
	Execute(ctx context.Context) error
}

type JobLocker interface {
	Lock(ctx context.Context, jobName string) (bool, error)
	Unlock(ctx context.Context, jobName string) error
}

// Runner runs a Job periodically with optional locking.
type Runner struct {
	job     Job
	locker  JobLocker
	jobName string
	period  time.Duration
	logger  *zap.Logger
}

// New returns a new CronJobRunner.
func New(job Job, locker JobLocker, jobName string, period time.Duration, logger *zap.Logger) *Runner {
	return &Runner{
		job:     job,
		locker:  locker,
		jobName: jobName,
		period:  period,
		logger:  logger,
	}
}

// Start runs the job every period until ctx is done.
func (r *Runner) Start(ctx context.Context) {
	ticker := time.NewTicker(r.period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		r.runOnce(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// runOnce tries to acquire the lock (if locker is set), runs the job, and releases the lock.
func (r *Runner) runOnce(ctx context.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			if r.logger != nil {
				r.logger.Error("panic in cron job: %v", zap.Any("recover", rec))
			}
		}
	}()

	locked, err := r.locker.Lock(ctx, r.jobName)
	if err != nil {
		if r.logger != nil {
			r.logger.Error("failed to acquire lock: %v", zap.Error(err))
		}
		return
	}
	if !locked {
		if r.logger != nil {
			r.logger.Info("job is already running, skipping execution")
		}
		return
	}
	defer func() {
		unlockErr := r.locker.Unlock(ctx, r.jobName)
		if unlockErr != nil && r.logger != nil {
			r.logger.Error("failed to release lock: %v", zap.Error(unlockErr))
		}
	}()

	err = r.job.Execute(ctx)
	if err != nil && r.logger != nil {
		r.logger.Error("failed to execute job: %v", zap.Error(err))
	}
}
