// This file is safe to edit. Once it exists it will not be overwritten

package api

import (
	"bytes"
	"io"

	"github.com/go-openapi/runtime/middleware"

	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/models"
	apiCalDav "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi/operations/cal_dav"
)

func (h *Handler) GetICalHandler(params apiCalDav.GetICalParams) middleware.Responder {
	ical, err := h.usecases.GetICal.Execute(params.HTTPRequest.Context(), params.Isu)
	if err != nil {
		return apiCalDav.NewGetICalInternalServerError().WithPayload(&models.Error{
			Error:   "InternalServerError",
			Message: err.Error(),
		})
	}
	if ical == nil {
		return apiCalDav.NewGetICalNotFound().WithPayload(&models.Error{
			Error:   "NotFound",
			Message: "iCal not found",
		})
	}

	return apiCalDav.NewGetICalOK().WithPayload(io.NopCloser(bytes.NewReader([]byte(ical.Serialize()))))
}
