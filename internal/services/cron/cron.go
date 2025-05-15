package cron

import (
	"context"

	"github.com/pkg/errors"
)

const (
	_batchSize = 10
)

type Service struct {
	client Client
}

func New(client Client) *Service {
	return &Service{
		client: client,
	}
}

// ScheduleSending splits ISUs into batches and sends them to the queue.
func (s *Service) ScheduleSending(ctx context.Context, isus []int64) error {
	for i := 0; i < len(isus); i += _batchSize {
		end := i + _batchSize
		if end > len(isus) {
			end = len(isus)
		}

		batch := isus[i:end]
		err := s.client.ScheduleSending(ctx, batch)
		if err != nil {
			return errors.Wrap(err, "schedule sending")
		}
	}

	return nil
}
