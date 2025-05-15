// This file is safe to edit. Once it exists it will not be overwritten

package api

import (
	"encoding/json"
	"net/http"

	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
)

func (h *Handler) SwaggerDocJSONHandler() http.Handler {
	specDoc, _ := loads.Analyzed(restapi.SwaggerJSON, "")

	b, _ := json.MarshalIndent(specDoc.Spec(), "", "  ")
	basePath := "/api/v1"
	handler := http.NotFoundHandler()

	return handlers.CORS()(middleware.Spec(basePath, b, handler))
}
