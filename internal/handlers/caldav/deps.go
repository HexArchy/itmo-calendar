package caldav

import (
	"context"

	ics "github.com/arran4/golang-ical"
)

// GetICal fetches a cached iCal calendar for the given ISU.
type GetICal interface {
	Execute(ctx context.Context, isu int64) (*ics.Calendar, error)
}

// Subscribe authenticates and creates a calendar subscription for the given ISU.
type Subscribe interface {
	Execute(ctx context.Context, isu int64, password string) error
}
