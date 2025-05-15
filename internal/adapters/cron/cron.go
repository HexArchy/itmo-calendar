package cron

import (
	"context"

	"github.com/hexarchy/itmo-calendar/pkg/rabbitmq"

	"github.com/pkg/errors"
)

type Adapter struct {
	client                   *rabbitmq.Client
	cronProcessScheduleQueue string
	sendScheduleQueue        string
}

func New(client *rabbitmq.Client, cronProcessScheduleQueue, sendScheduleQueue string) *Adapter {
	return &Adapter{
		client:                   client,
		cronProcessScheduleQueue: cronProcessScheduleQueue,
		sendScheduleQueue:        sendScheduleQueue,
	}
}

func (a *Adapter) SendCronTask(ctx context.Context) error {
	msg, err := rabbitmq.NewMessage(struct{}{}, nil)
	if err != nil {
		return errors.Wrap(err, "new message")
	}

	err = a.client.SendMessage(ctx, a.cronProcessScheduleQueue, msg)
	if err != nil {
		return errors.Wrap(err, "send cron task")
	}

	return nil
}

func (a *Adapter) ScheduleSending(ctx context.Context, isus []int64) error {
	payload := struct {
		ISUs []int64 `json:"isus"`
	}{
		ISUs: isus,
	}

	msg, err := rabbitmq.NewMessage(payload, nil)
	if err != nil {
		return errors.Wrap(err, "new message")
	}

	err = a.client.SendMessage(ctx, a.sendScheduleQueue, msg)
	if err != nil {
		return errors.Wrap(err, "send schedule")
	}

	return nil
}
