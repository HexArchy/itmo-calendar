package preparesendschedule

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const _batchSize = 100

type UseCase struct {
	cron   Cron
	users  Users
	logger *zap.Logger
}

func New(cron Cron, users Users, logger *zap.Logger) *UseCase {
	return &UseCase{
		cron:   cron,
		users:  users,
		logger: logger,
	}
}

func (u *UseCase) Execute(ctx context.Context) error {
	offset := 0

	for {
		users, err := u.users.GetBatch(ctx, _batchSize, offset)
		if err != nil {
			return errors.Wrap(err, "get batch users")
		}

		if len(users) == 0 {
			break
		}

		isus := make([]int64, 0, len(users))
		for _, user := range users {
			isus = append(isus, user.ISU)
		}

		err = u.cron.ScheduleSending(ctx, isus)
		if err != nil {
			return errors.Wrap(err, "schedule sending")
		}

		offset += _batchSize

		if len(users) < _batchSize {
			break
		}
	}

	return nil
}
