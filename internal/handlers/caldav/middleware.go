package caldav

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

// AuthMiddleware handles HTTP Basic Auth for CalDAV requests.
// It parses the ISU from the username, auto-subscribes if no calendar exists,
// and stores the ISU and password in the request context.
type AuthMiddleware struct {
	getICal   GetICal
	subscribe Subscribe
	logger    *zap.Logger
}

func NewAuthMiddleware(getICal GetICal, subscribe Subscribe, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		getICal:   getICal,
		subscribe: subscribe,
		logger:    logger.With(zap.String("component", "caldav_auth")),
	}
}

func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			m.unauthorized(w)
			return
		}

		isu, err := strconv.ParseInt(username, 10, 64)
		if err != nil {
			m.unauthorized(w)
			return
		}

		ctx := r.Context()

		// Fast path: check cached calendar exists.
		cal, err := m.getICal.Execute(ctx, isu)
		if err != nil && !errors.Is(err, entities.ErrNotFound) {
			m.logger.Error("get ical", zap.Error(err), zap.Int64("isu", isu))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if cal == nil {
			// No cached data — auto-subscribe.
			if err := m.subscribe.Execute(ctx, isu, password); err != nil {
				m.logger.Warn("subscribe failed", zap.Error(err), zap.Int64("isu", isu))
				m.unauthorized(w)
				return
			}
		}

		ctx = ContextWithISU(ctx, isu)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="ITMO Calendar"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
