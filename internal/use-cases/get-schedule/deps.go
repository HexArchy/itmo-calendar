package getschedule

import (
	"context"

	ics "github.com/arran4/golang-ical"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

type CalDav interface {
	Get(ctx context.Context, isu int64) (entities.CalDav, error)
}

type ICal interface {
	Parse(_ context.Context, cal *ics.Calendar) ([]entities.DaySchedule, error)
}
