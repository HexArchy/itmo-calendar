package caldav

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	ical "github.com/emersion/go-ical"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"go.uber.org/zap"
)

var (
	errNoISU    = errors.New("no ISU in context")
	errReadOnly = errors.New("read-only calendar")
)

// Backend implements caldav.Backend as a read-only CalDAV server.
type Backend struct {
	getICal GetICal
	logger  *zap.Logger
}

// NewBackend creates a new CalDAV backend.
func NewBackend(getICal GetICal, logger *zap.Logger) *Backend {
	return &Backend{
		getICal: getICal,
		logger:  logger.With(zap.String("component", "caldav_backend")),
	}
}

// CurrentUserPrincipal returns the principal URL for the authenticated user.
func (b *Backend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	isu, ok := ISUFromContext(ctx)
	if !ok {
		return "", webdav.NewHTTPError(http.StatusUnauthorized, errNoISU)
	}
	return fmt.Sprintf("/caldav/%d/", isu), nil
}

// CalendarHomeSetPath returns the calendar home set URL.
func (b *Backend) CalendarHomeSetPath(ctx context.Context) (string, error) {
	isu, ok := ISUFromContext(ctx)
	if !ok {
		return "", webdav.NewHTTPError(http.StatusUnauthorized, errNoISU)
	}
	return fmt.Sprintf("/caldav/%d/calendars/", isu), nil
}

// ListCalendars returns the single schedule calendar.
func (b *Backend) ListCalendars(ctx context.Context) ([]caldav.Calendar, error) {
	isu, ok := ISUFromContext(ctx)
	if !ok {
		return nil, webdav.NewHTTPError(http.StatusUnauthorized, errNoISU)
	}
	return []caldav.Calendar{b.makeCalendar(isu)}, nil
}

// GetCalendar returns calendar metadata.
func (b *Backend) GetCalendar(ctx context.Context, _ string) (*caldav.Calendar, error) {
	isu, ok := ISUFromContext(ctx)
	if !ok {
		return nil, webdav.NewHTTPError(http.StatusUnauthorized, errNoISU)
	}
	cal := b.makeCalendar(isu)
	return &cal, nil
}

// CreateCalendar is forbidden (read-only).
func (b *Backend) CreateCalendar(_ context.Context, _ *caldav.Calendar) error {
	return webdav.NewHTTPError(http.StatusForbidden, errReadOnly)
}

// ListCalendarObjects returns all events as individual calendar objects.
func (b *Backend) ListCalendarObjects(
	ctx context.Context,
	_ string,
	_ *caldav.CalendarCompRequest,
) ([]caldav.CalendarObject, error) {
	return b.loadObjects(ctx)
}

// GetCalendarObject returns a single calendar object by path.
func (b *Backend) GetCalendarObject(
	ctx context.Context,
	path string,
	_ *caldav.CalendarCompRequest,
) (*caldav.CalendarObject, error) {
	uid := uidFromPath(path)
	objects, err := b.loadObjects(ctx)
	if err != nil {
		return nil, err
	}
	for i := range objects {
		if uidFromPath(objects[i].Path) == uid {
			return &objects[i], nil
		}
	}
	return nil, webdav.NewHTTPError(
		http.StatusNotFound,
		fmt.Errorf("calendar object %q not found", uid),
	)
}

// QueryCalendarObjects returns objects matching the query filter.
func (b *Backend) QueryCalendarObjects(
	ctx context.Context,
	_ string,
	query *caldav.CalendarQuery,
) ([]caldav.CalendarObject, error) {
	objects, err := b.loadObjects(ctx)
	if err != nil {
		return nil, err
	}

	// If the query has a time-range filter on VEVENT, apply it.
	start, end := extractTimeRange(query)
	if start.IsZero() && end.IsZero() {
		return objects, nil
	}

	var filtered []caldav.CalendarObject
	for _, obj := range objects {
		if matchesTimeRange(obj.Data, start, end) {
			filtered = append(filtered, obj)
		}
	}
	return filtered, nil
}

// PutCalendarObject is forbidden (read-only).
func (b *Backend) PutCalendarObject(
	_ context.Context,
	_ string,
	_ *ical.Calendar,
	_ *caldav.PutCalendarObjectOptions,
) (*caldav.CalendarObject, error) {
	return nil, webdav.NewHTTPError(http.StatusForbidden, errReadOnly)
}

// DeleteCalendarObject is forbidden (read-only).
func (b *Backend) DeleteCalendarObject(_ context.Context, _ string) error {
	return webdav.NewHTTPError(http.StatusForbidden, errReadOnly)
}

func (b *Backend) makeCalendar(isu int64) caldav.Calendar {
	return caldav.Calendar{
		Path:                  fmt.Sprintf("/caldav/%d/calendars/schedule/", isu),
		Name:                  "ITMO Schedule",
		Description:           "ITMO University class schedule",
		SupportedComponentSet: []string{ical.CompEvent},
	}
}

func (b *Backend) loadObjects(ctx context.Context) ([]caldav.CalendarObject, error) {
	isu, ok := ISUFromContext(ctx)
	if !ok {
		return nil, webdav.NewHTTPError(http.StatusUnauthorized, errNoISU)
	}

	cal, err := b.getICal.Execute(ctx, isu)
	if err != nil {
		return nil, fmt.Errorf("get ical: %w", err)
	}

	events := cal.Events()
	objects := make([]caldav.CalendarObject, 0, len(events))

	for _, event := range events {
		uid := event.Id()
		if uid == "" {
			continue
		}

		serialized := wrapSingleEvent(event)
		emCal, convErr := convertToEmersionCal(serialized)
		if convErr != nil {
			b.logger.Warn("convert event", zap.Error(convErr), zap.String("uid", uid))
			continue
		}

		objects = append(objects, caldav.CalendarObject{
			Path: fmt.Sprintf("/caldav/%d/calendars/schedule/%s.ics", isu, uid),
			ETag: computeETag(serialized),
			Data: emCal,
		})
	}

	return objects, nil
}
