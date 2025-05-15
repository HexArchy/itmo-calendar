package schedules

import (
	"context"
	"time"

	"github.com/hexarchy/itmo-calendar/internal/entities"

	"github.com/pkg/errors"
)

// Service provides schedule-related operations.
type Service struct {
	schedule   ScheduleRepo
	tokens     Tokens
	userTokens UserTokensRepo
}

// New creates a new Schedule.
func New(schedule ScheduleRepo, tokens Tokens, userTokens UserTokensRepo) *Service {
	return &Service{
		schedule:   schedule,
		tokens:     tokens,
		userTokens: userTokens,
	}
}

// GetByCreds retrieves schedule using ISU credentials and stores tokens.
func (s *Service) GetByCreds(ctx context.Context, isu int64, password string, from, to time.Time) ([]entities.DaySchedule, error) {
	tokens, err := s.tokens.Get(ctx, isu, password)
	if err != nil {
		return nil, errors.Wrap(err, "get tokens")
	}

	err = s.userTokens.UpsertUserTokens(ctx, tokens)
	if err != nil {
		return nil, errors.Wrap(err, "upsert tokens")
	}

	schedule, err := s.schedule.Get(ctx, tokens.AccessToken, from, to)
	if err != nil {
		return nil, errors.Wrap(err, "get schedule")
	}

	return schedule, nil
}

// GetByISU retrieves schedule for a user, refreshing tokens if needed.
func (s *Service) GetByISU(ctx context.Context, isu int64, from, to time.Time) ([]entities.DaySchedule, error) {
	tokens, err := s.userTokens.Get(ctx, isu)
	if err != nil {
		return nil, errors.Wrap(err, "get user tokens")
	}

	if tokens == nil {
		return nil, errors.New("user tokens not found")
	}

	now := time.Now()
	if now.After(tokens.AccessTokenExpiresAt) {
		if now.After(tokens.RefreshTokenExpiresAt) {
			return nil, errors.New("refresh token expired")
		}

		newTokens, err := s.tokens.Refresh(ctx, isu, tokens.RefreshToken)
		if err != nil {
			return nil, errors.Wrap(err, "refresh tokens")
		}

		err = s.userTokens.UpsertUserTokens(ctx, newTokens)
		if err != nil {
			return nil, errors.Wrap(err, "upsert refreshed tokens")
		}

		tokens = newTokens
	}

	schedule, err := s.schedule.Get(ctx, tokens.AccessToken, from, to)
	if err != nil {
		return nil, errors.Wrap(err, "get schedule")
	}

	return schedule, nil
}
