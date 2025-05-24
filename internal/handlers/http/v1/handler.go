// This file is safe to edit. Once it exists it will not be overwritten

package api

import (
	"net/http"

	"github.com/go-openapi/loads"
	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/app/container"
	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi"
	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi/operations"

	apiCalDav "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi/operations/cal_dav"
	apiSchedule "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi/operations/schedule"
	apiSystem "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi/operations/system"
)

type Handler struct {
	ops      *operations.ItmoCalendarAPI
	usecases *container.UseCases
	logger   *zap.Logger
}

func NewHandler(usecases *container.UseCases, logger *zap.Logger) (*Handler, error) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}
	r := &Handler{
		ops:      operations.NewItmoCalendarAPI(swaggerSpec),
		usecases: usecases,
		logger:   logger.With(zap.String("component", "api_handler")),
	}
	r.setUpHandlers()

	return r, nil
}

func (h *Handler) handlerFor(method, path string) http.Handler {
	r, _ := h.ops.HandlerFor(method, path)

	return r
}

func (h *Handler) setUpHandlers() {
	h.ops.SystemHealthCheckHandler = apiSystem.HealthCheckHandlerFunc(h.HealthCheckHandler)
	h.ops.CalDavGetICalHandler = apiCalDav.GetICalHandlerFunc(h.GetICalHandler)
	h.ops.CalDavSubscribeScheduleHandler = apiCalDav.SubscribeScheduleHandlerFunc(h.SubscribeScheduleHandler)
	h.ops.ScheduleGetScheduleHandler = apiSchedule.GetScheduleHandlerFunc(h.GetScheduleHandler)

	// You can add your middleware to concrete route
	// h.ops.AddMiddlewareFor("%method%", "%route%", %middlewareBuilder%)

	// You can add your global middleware
	// h.ops.AddGlobalMiddleware(%middlewareBuilder%)

	configureAPI(h.ops)
}
