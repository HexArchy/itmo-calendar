// This file is safe to edit. Once it exists it will not be overwritten

package api

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/models"
	apiCalDav "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi/operations/cal_dav"
)

func (h *Handler) SubscribeScheduleHandler(params apiCalDav.SubscribeScheduleParams) middleware.Responder {
	if params.Body.Isu == nil || params.Body.Password == nil {
		return apiCalDav.NewSubscribeScheduleBadRequest().WithPayload(&models.Error{
			Error:   "BadRequest",
			Message: "ISU and password are required",
		})
	}

	err := h.usecases.SubscirbeSchedule.Execute(params.HTTPRequest.Context(), *params.Body.Isu, *params.Body.Password)
	if err != nil {
		return apiCalDav.NewSubscribeScheduleInternalServerError().WithPayload(&models.Error{
			Error:   "InternalServerError",
			Message: err.Error(),
		})
	}

	return apiCalDav.NewSubscribeScheduleOK().WithPayload(&models.SubscribeResponse{
		Message: "Subscription successful. iCal generated.",
	})
}
