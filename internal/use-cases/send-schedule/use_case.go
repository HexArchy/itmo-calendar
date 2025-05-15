package sendschedule

import (
	"context"
	"time"

	"github.com/hexarchy/itmo-calendar/internal/entities"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	_fromDays      = 30  // days.
	_toDays        = 120 // days.
	_maxConcurrent = 5
)

type UseCase struct {
	schedules Schedules
	users     Users
	iCal      ICal
	calDav    CalDav
	logger    *zap.Logger
}

func New(schedules Schedules, users Users, iCal ICal, calDav CalDav, logger *zap.Logger) *UseCase {
	return &UseCase{
		schedules: schedules,
		users:     users,
		iCal:      iCal,
		calDav:    calDav,
		logger:    logger,
	}
}

func (u *UseCase) Execute(ctx context.Context, isus []int64) error {
	users, err := u.users.FindByIDs(ctx, isus)
	if err != nil {
		return errors.Wrap(err, "find by ids")
	}

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(_maxConcurrent)

	for _, user := range users {
		g.Go(func() error {
			if errProcess := u.processSending(gctx, user); errProcess != nil {
				u.logger.Error("failed to process sending", zap.Error(errProcess), zap.Int64("isu", user.ISU))
				return nil
			}
			u.logger.Debug("schedule sent", zap.Int64("isu", user.ISU))
			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		return errors.Wrap(err, "process users")
	}

	return nil
}

func (u *UseCase) processSending(ctx context.Context, user entities.User) error {
	from := time.Now().AddDate(0, 0, -_fromDays)
	to := time.Now().AddDate(0, 0, _toDays)

	schedule, err := u.schedules.GetByISU(ctx, user.ISU, from, to)
	if err != nil {
		return errors.Wrap(err, "get schedule")
	}

	ical, err := u.iCal.Generate(ctx, schedule)
	if err != nil {
		return errors.Wrap(err, "generate iCal")
	}

	err = u.calDav.Create(ctx, user, ical)
	if err != nil {
		return errors.Wrap(err, "send schedule")
	}

	return nil
}
