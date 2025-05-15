package preparesendschedule

import (
	"context"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

type Cron interface {
	ScheduleSending(ctx context.Context, isus []int64) error
}

type Users interface {
	GetAll(ctx context.Context) ([]entities.User, error)
}
