package getical

import (
	"context"

	ics "github.com/arran4/golang-ical"
	"github.com/pkg/errors"
)

type UseCase struct {
	calDav CalDav
}

func New(calDav CalDav) *UseCase {
	return &UseCase{
		calDav: calDav,
	}
}

func (u *UseCase) Execute(ctx context.Context, isu int64) (*ics.Calendar, error) {
	calDav, err := u.calDav.Get(ctx, isu)
	if err != nil {
		return nil, errors.Wrap(err, "get caldav")
	}

	return calDav.ICal, nil
}
