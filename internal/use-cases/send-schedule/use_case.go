package sendschedule

import (
	"context"
	"time"

	"github.com/hexarchy/itmo-calendar/internal/entities"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_defaultFromTimePeriod = 31  // days.
	_defaultToTimePeriod   = 120 // days.
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

	for _, user := range users {
		if errProcess := u.processSending(ctx, user); errProcess != nil {
			u.logger.Error("failed to process sending", zap.Error(errProcess), zap.Int64("isu", user.ISU))
			continue
		}
		u.logger.Debug("schedule sent", zap.Int64("isu", user.ISU))
	}

	return nil
}

func (u *UseCase) processSending(ctx context.Context, user entities.User) error {
	from := time.Now().AddDate(0, 0, -_defaultFromTimePeriod)
	to := time.Now().AddDate(0, 0, _defaultToTimePeriod)

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
