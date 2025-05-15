package sendschedule

import (
	"context"
	"encoding/json"

	"github.com/hexarchy/itmo-calendar/pkg/rabbitmq"

	"go.uber.org/zap"
)

type UseCase interface {
	Execute(ctx context.Context, isus []int64) error
}

// Worker handles tasks from the send-schedule queue.
type Worker struct {
	rabbit  *rabbitmq.Client
	useCase UseCase
	queue   string
	logger  *zap.Logger
}

// New returns a new Worker.
func New(rabbit *rabbitmq.Client, useCase UseCase, queue string, logger *zap.Logger) *Worker {
	return &Worker{
		rabbit:  rabbit,
		useCase: useCase,
		queue:   queue,
		logger:  logger.With(zap.String("worker", "send-schedule")),
	}
}

// Start launches the worker to consume tasks.
func (w *Worker) Start(ctx context.Context) error {
	processFunc := func(ctx context.Context, msg *rabbitmq.Message) error {
		var payload struct {
			ISUs []int64 `json:"isus"`
		}

		err := json.Unmarshal(msg.Body, &payload)
		if err != nil {
			w.logger.Error("failed to unmarshal payload", zap.Error(err))
			return err
		}

		w.logger.Debug("processing send-schedule task", zap.String("queue", w.queue), zap.Any("payload", payload))

		err = w.useCase.Execute(ctx, payload.ISUs)
		if err != nil {
			w.logger.Error("failed to execute send-schedule use case", zap.Error(err))
			return err
		}

		return nil
	}

	err := w.rabbit.DefineQueue(
		ctx,
		w.queue,
		4, // TODO: make it configurable.
		4, // TODO: make it configurable.
		processFunc,
	)
	if err != nil {
		w.logger.Error("failed to define queue", zap.Error(err))
		return err
	}

	w.logger.Info("send-schedule worker started")
	<-ctx.Done()

	return nil
}
