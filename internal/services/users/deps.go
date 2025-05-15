package users

import (
	"context"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

type Repository interface {
	GetAll(ctx context.Context) ([]entities.User, error)
	FindByIDs(ctx context.Context, isus []int64) ([]entities.User, error)
	Create(ctx context.Context, isu int64) (*entities.User, error)
}
