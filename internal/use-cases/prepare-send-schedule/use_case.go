package preparesendschedule

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

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
	users, err := u.users.GetAll(ctx)
	if err != nil {
		return errors.Wrap(err, "get all users")
	}

	isus := make([]int64, 0, len(users))
	for _, user := range users {
		isus = append(isus, user.ISU)
	}

	err = u.cron.ScheduleSending(ctx, isus)
	if err != nil {
		return errors.Wrap(err, "schedule sending")
	}

	return nil
}
