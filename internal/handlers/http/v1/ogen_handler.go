package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/entities"
	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/gen"

	getical "github.com/hexarchy/itmo-calendar/internal/use-cases/get-ical"
	getschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/get-schedule"
	subscribeschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/subscribe-schedule"
)

// OgenHandler implements gen.Handler using application use-cases.
type OgenHandler struct {
	getICal     *getical.UseCase
	getSchedule *getschedule.UseCase
	subscribe   *subscribeschedule.UseCase
	logger      *zap.Logger
}

// NewOgenHandler creates a new OgenHandler.
func NewOgenHandler(
	getICal *getical.UseCase,
	getSchedule *getschedule.UseCase,
	subscribe *subscribeschedule.UseCase,
	logger *zap.Logger,
) *OgenHandler {
	return &OgenHandler{
		getICal:     getICal,
		getSchedule: getSchedule,
		subscribe:   subscribe,
		logger:      logger.With(zap.String("component", "ogen_handler")),
	}
}

// HealthCheck implements gen.Handler.
func (h *OgenHandler) HealthCheck(_ context.Context) error {
	return nil
}

// GetICal implements gen.Handler.
func (h *OgenHandler) GetICal(ctx context.Context, params gen.GetICalParams) (gen.GetICalOK, error) {
	cal, err := h.getICal.Execute(ctx, params.Isu)
	if err != nil {
		h.logger.Error("get ical failed", zap.Error(err), zap.Int64("isu", params.Isu))
		return gen.GetICalOK{}, err
	}

	return gen.GetICalOK{Data: strings.NewReader(cal.Serialize())}, nil
}

// GetSchedule implements gen.Handler.
func (h *OgenHandler) GetSchedule(ctx context.Context, params gen.GetScheduleParams) ([]gen.ScheduleItem, error) {
	schedule, err := h.getSchedule.Execute(ctx, params.Isu)
	if err != nil {
		h.logger.Error("get schedule failed", zap.Error(err), zap.Int64("isu", params.Isu))
		return nil, err
	}

	items := make([]gen.ScheduleItem, 0, len(schedule))
	for _, day := range schedule {
		lessons := make([]gen.Lesson, 0, len(day.Lessons))
		for _, l := range day.Lessons {
			lesson := gen.Lesson{
				Subject:     l.Subject,
				Type:        l.Type,
				TeacherName: l.TeacherName,
				Room:        l.Room,
				Building:    l.Building,
				Format:      l.Format,
				Group:       l.Group,
				TimeStart:   l.Start,
				TimeEnd:     l.End,
			}
			if l.Note != "" {
				lesson.Note = gen.NewOptString(l.Note)
			}
			if l.ZoomURL != "" {
				lesson.ZoomURL = gen.NewOptString(l.ZoomURL)
			}
			lessons = append(lessons, lesson)
		}
		items = append(items, gen.ScheduleItem{
			Date:    day.Date,
			Lessons: lessons,
		})
	}

	return items, nil
}

// SubscribeSchedule implements gen.Handler.
func (h *OgenHandler) SubscribeSchedule(
	ctx context.Context,
	req *gen.SubscribeRequest,
) (*gen.SubscribeResponse, error) {
	err := h.subscribe.Execute(ctx, req.Isu, req.Password)
	if err != nil {
		h.logger.Error("subscribe failed", zap.Error(err), zap.Int64("isu", req.Isu))
		return nil, err
	}

	return &gen.SubscribeResponse{
		Message: gen.NewOptString("Subscription successful. iCal generated."),
	}, nil
}

// NewError implements gen.Handler.
func (h *OgenHandler) NewError(_ context.Context, err error) *gen.ErrorStatusCode {
	code := http.StatusInternalServerError
	title := "Internal Server Error"

	if errors.Is(err, entities.ErrNotFound) {
		code = http.StatusNotFound
		title = "Not Found"
	}

	return &gen.ErrorStatusCode{
		StatusCode: code,
		Response: gen.Error{
			Error:   gen.NewOptString(title),
			Message: gen.NewOptString(err.Error()),
		},
	}
}
