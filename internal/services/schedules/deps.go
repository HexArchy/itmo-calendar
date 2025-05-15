package schedules

import (
	"context"
	"time"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

type UserTokensRepo interface {
	Get(ctx context.Context, isu int64) (*entities.UserTokens, error)
	UpsertUserTokens(ctx context.Context, tokens *entities.UserTokens) error
}

type Tokens interface {
	Get(ctx context.Context, isu int64, password string) (*entities.UserTokens, error)
	Refresh(ctx context.Context, isu int64, refreshToken string) (*entities.UserTokens, error)
}

type ScheduleRepo interface {
	Get(ctx context.Context, token string, from, to time.Time) ([]entities.DaySchedule, error)
}
