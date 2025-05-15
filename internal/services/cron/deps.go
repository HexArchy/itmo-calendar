package cron

import (
	"context"
)

type Client interface {
	ScheduleSending(ctx context.Context, isus []int64) error
	SendCronTask(ctx context.Context) error
}
