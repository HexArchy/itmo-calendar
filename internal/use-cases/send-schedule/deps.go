package sendschedule

import (
	"context"
	"time"

	ics "github.com/arran4/golang-ical"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

type Schedules interface {
	GetByISU(ctx context.Context, isu int64, from, to time.Time) ([]entities.DaySchedule, error)
}

type Users interface {
	FindByIDs(ctx context.Context, isus []int64) ([]entities.User, error)
}

type ICal interface {
	Generate(ctx context.Context, schedule []entities.DaySchedule) (*ics.Calendar, error)
}

type CalDav interface {
	Create(ctx context.Context, user entities.User, ical *ics.Calendar) error
}
