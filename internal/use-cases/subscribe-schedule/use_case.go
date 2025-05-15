package subscribeschedule

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const _period = 120 // days

type UseCase struct {
	schedules Schedules
	users     Users
	iCal      ICal
	caldav    CalDav
	logger    *zap.Logger
}

func New(schedules Schedules, users Users, iCal ICal, caldav CalDav, logger *zap.Logger) *UseCase {
	return &UseCase{
		schedules: schedules,
		users:     users,
		iCal:      iCal,
		caldav:    caldav,
		logger:    logger,
	}
}
func (u *UseCase) Execute(ctx context.Context, isu int64, password string) error {
	from := time.Now().AddDate(0, 0, -30)
	to := time.Now().AddDate(0, 0, _period)

	schedule, err := u.schedules.GetByCreds(ctx, isu, password, from, to)
	if err != nil {
		return errors.Wrap(err, "get schedule")
	}

	user, err := u.users.Create(ctx, isu)
	if err != nil {
		return errors.Wrap(err, "create user")
	}

	ical, err := u.iCal.Generate(ctx, schedule)
	if err != nil {
		return errors.Wrap(err, "generate iCal")
	}

	err = u.caldav.Create(ctx, *user, ical)
	if err != nil {
		return errors.Wrap(err, "send schedule")
	}

	return nil
}
