// This file is safe to edit. Once it exists it will not be overwritten

package api

import (
	"github.com/go-openapi/runtime/middleware"

	apiSystem "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi/operations/system"
)

func (h *Handler) HealthCheckHandler(_ apiSystem.HealthCheckParams) middleware.Responder {
	return apiSystem.NewHealthCheckOK()
}
