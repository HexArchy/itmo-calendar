// This file is safe to edit. Once it exists it will not be overwritten

package api

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

func (h *Handler) SwaggerDocUIHandler() http.Handler {
	specDoc, _ := loads.Analyzed(restapi.SwaggerJSON, "")

	b, _ := json.MarshalIndent(specDoc.Spec(), "", "  ")

	basePath := "/api/v1"
	handler := http.NotFoundHandler()

	swaggerUIOpts := middleware.SwaggerUIOpts{
		BasePath: basePath,
		Title:    "Itmo Calendar",
		SpecURL:  path.Join(basePath, "/swagger.json"),
	}

	return middleware.Spec(basePath, b, middleware.SwaggerUI(swaggerUIOpts, handler))
}
