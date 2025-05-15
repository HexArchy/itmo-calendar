package caldav

import (
	"context"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

type Repo interface {
	Create(ctx context.Context, caldav entities.CalDav) error
	Get(ctx context.Context, isu int64) (entities.CalDav, error)
}
