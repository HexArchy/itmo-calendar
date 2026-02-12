package http

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/entities"
	getical "github.com/hexarchy/itmo-calendar/internal/use-cases/get-ical"
	subscribeschedule "github.com/hexarchy/itmo-calendar/internal/use-cases/subscribe-schedule"
)

// CalHandler serves iCal feed with HTTP Basic Auth for Apple Calendar subscriptions.
type CalHandler struct {
	getICal   *getical.UseCase
	subscribe *subscribeschedule.UseCase
	logger    *zap.Logger
}

func NewCalHandler(getICal *getical.UseCase, subscribe *subscribeschedule.UseCase, logger *zap.Logger) *CalHandler {
	return &CalHandler{
		getICal:   getICal,
		subscribe: subscribe,
		logger:    logger.With(zap.String("handler", "cal")),
	}
}

func (h *CalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="ITMO Calendar"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isu, err := strconv.ParseInt(username, 10, 64)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="ITMO Calendar"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	// Fast path: return cached iCal.
	cal, err := h.getICal.Execute(ctx, isu)
	if err != nil && !errors.Is(err, entities.ErrNotFound) {
		h.logger.Error("get ical", zap.Error(err), zap.Int64("isu", isu))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if cal == nil {
		// No cached data — subscribe (ITMO OAuth + fetch + generate).
		if err := h.subscribe.Execute(ctx, isu, password); err != nil {
			h.logger.Warn("subscribe failed", zap.Error(err), zap.Int64("isu", isu))
			w.Header().Set("WWW-Authenticate", `Basic realm="ITMO Calendar"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		cal, err = h.getICal.Execute(ctx, isu)
		if err != nil {
			h.logger.Error("get ical after subscribe", zap.Error(err), zap.Int64("isu", isu))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := cal.SerializeTo(w); err != nil {
		h.logger.Error("serialize ical", zap.Error(err), zap.Int64("isu", isu))
	}
}
