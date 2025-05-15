package getical

import (
	"context"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

type CalDav interface {
	Get(ctx context.Context, isu int64) (entities.CalDav, error)
}
