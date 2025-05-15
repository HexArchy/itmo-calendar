package users

import (
	"context"

	"github.com/hexarchy/itmo-calendar/internal/entities"

	"github.com/pkg/errors"
)

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, isu int64) (*entities.User, error) {
	createdUser, err := s.repo.Create(ctx, isu)
	if err != nil {
		return nil, errors.Wrap(err, "create user")
	}

	return createdUser, nil
}

func (s *Service) GetAll(ctx context.Context) ([]entities.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get all users")
	}
	return users, nil
}

func (s *Service) FindByIDs(ctx context.Context, isus []int64) ([]entities.User, error) {
	users, err := s.repo.FindByIDs(ctx, isus)
	if err != nil {
		return nil, errors.Wrap(err, "find users by ids")
	}
	return users, nil
}
