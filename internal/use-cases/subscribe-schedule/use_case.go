package subscribeschedule

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_fromDays = 30  // days.
	_toDays   = 120 // days.
)

type UseCase struct {
	tx        TxRunner
	schedules Schedules
	users     Users
	iCal      ICal
	caldav    CalDav
	logger    *zap.Logger
}

func New(tx TxRunner, schedules Schedules, users Users, iCal ICal, caldav CalDav, logger *zap.Logger) *UseCase {
	return &UseCase{
		tx:        tx,
		schedules: schedules,
		users:     users,
		iCal:      iCal,
		caldav:    caldav,
		logger:    logger,
	}
}

func (u *UseCase) Execute(ctx context.Context, isu int64, password string) error {
	return u.tx.Run(ctx, func(ctx context.Context) error {
		user, err := u.users.Create(ctx, isu)
		if err != nil {
			return errors.Wrap(err, "create user")
		}

		from := time.Now().AddDate(0, 0, -_fromDays)
		to := time.Now().AddDate(0, 0, _toDays)

		schedule, err := u.schedules.GetByCreds(ctx, isu, password, from, to)
		if err != nil {
			return errors.Wrap(err, "get schedule")
		}

		ical, err := u.iCal.Generate(ctx, schedule)
		if err != nil {
			return errors.Wrap(err, "generate iCal")
		}

		if err = u.caldav.Create(ctx, *user, ical); err != nil {
			return errors.Wrap(err, "send schedule")
		}

		return nil
	})
}
