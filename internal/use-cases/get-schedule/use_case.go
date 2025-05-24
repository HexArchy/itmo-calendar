package getschedule

import (
	"context"

	"github.com/hexarchy/itmo-calendar/internal/entities"
	"github.com/pkg/errors"
)

type UseCase struct {
	calDav CalDav
	ical   ICal
}

func New(calDav CalDav, ical ICal) *UseCase {
	return &UseCase{
		calDav: calDav,
		ical:   ical,
	}
}

func (u *UseCase) Execute(ctx context.Context, isu int64) ([]entities.DaySchedule, error) {
	caldav, err := u.calDav.Get(ctx, isu)
	if err != nil {
		return nil, errors.Wrap(err, "get caldav")
	}

	schedule, err := u.ical.Parse(ctx, caldav.ICal)
	if err != nil {
		return nil, errors.Wrap(err, "parse calendar")
	}

	return schedule, nil
}
