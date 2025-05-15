package caldav

import (
	"context"

	ics "github.com/arran4/golang-ical"
	"github.com/pkg/errors"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

type Service struct {
	repo Repo
}

func New(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, user entities.User, ical *ics.Calendar) error {
	err := s.repo.Create(ctx, entities.CalDav{
		ICal: ical,
		ISU:  user.ISU,
	})
	if err != nil {
		return errors.Wrap(err, "create caldav")
	}

	return nil
}

func (s *Service) Get(ctx context.Context, isu int64) (entities.CalDav, error) {
	calDav, err := s.repo.Get(ctx, isu)
	if err != nil {
		return entities.CalDav{}, errors.Wrap(err, "get caldav")
	}

	return calDav, nil
}
