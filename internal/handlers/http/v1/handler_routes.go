// Code generated; DO NOT EDIT.

package api

import (
	"fmt"
	"strings"

	"github.com/gorilla/mux"
)

const version = "1.0.0"

func (h *Handler) AddRoutes(router *mux.Router) {

	router.Handle("/{isu}/ical", h.handlerFor("GET", "/{isu}/ical")).Methods("GET")
	router.Handle("/{isu}/schedule", h.handlerFor("GET", "/{isu}/schedule")).Methods("GET")
	router.Handle("/health", h.handlerFor("GET", "/health")).Methods("GET")
	router.Handle("/subscribe", h.handlerFor("POST", "/subscribe")).Methods("POST")

	router.Handle("/swagger.json", h.SwaggerDocJSONHandler()).Methods("GET")
	router.Handle("/docs", h.SwaggerDocUIHandler()).Methods("GET")
}

func (h *Handler) GetVersion() string {
	return fmt.Sprintf("v%s", strings.Split(version, ".")[0])
}
