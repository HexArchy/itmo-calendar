package caldav

import (
	"crypto/sha256"
	"fmt"
	"path"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	ical "github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

// wrapSingleEvent creates a minimal VCALENDAR containing a single VEVENT.
func wrapSingleEvent(event *ics.VEvent) string {
	cal := ics.NewCalendar()
	cal.AddVEvent(event)
	return cal.Serialize()
}

// convertToEmersionCal parses an arran4/golang-ical serialized string
// into an emersion/go-ical Calendar.
func convertToEmersionCal(serialized string) (*ical.Calendar, error) {
	dec := ical.NewDecoder(strings.NewReader(serialized))
	return dec.Decode()
}

// computeETag returns a quoted ETag based on SHA-256 of the serialized event.
func computeETag(serialized string) string {
	h := sha256.Sum256([]byte(serialized))
	return fmt.Sprintf(`"%x"`, h[:8])
}

// uidFromPath extracts the UID from a CalDAV object path.
// Example: /caldav/123/calendars/schedule/abc.ics → abc.
func uidFromPath(p string) string {
	base := path.Base(p)
	return strings.TrimSuffix(base, ".ics")
}

// extractTimeRange extracts the time-range from a CalendarQuery's VEVENT comp filter.
func extractTimeRange(query *caldav.CalendarQuery) (time.Time, time.Time) {
	if query == nil {
		return time.Time{}, time.Time{}
	}
	// Look for VEVENT comp filter inside VCALENDAR.
	for _, comp := range query.CompFilter.Comps {
		if comp.Name == ical.CompEvent {
			return comp.Start, comp.End
		}
	}
	// Fallback: check the top-level filter itself.
	return query.CompFilter.Start, query.CompFilter.End
}

// matchesTimeRange checks if any VEVENT in the calendar overlaps [start, end).
func matchesTimeRange(cal *ical.Calendar, start, end time.Time) bool {
	if cal == nil {
		return false
	}
	for _, child := range cal.Children {
		if child.Name != ical.CompEvent {
			continue
		}
		event := ical.Event{Component: child}
		dtStart, err := event.DateTimeStart(nil)
		if err != nil || dtStart.IsZero() {
			continue
		}
		dtEnd, err := event.DateTimeEnd(nil)
		if err != nil || dtEnd.IsZero() {
			dtEnd = dtStart.Add(time.Hour)
		}

		if !end.IsZero() && dtStart.After(end) {
			continue
		}
		if !start.IsZero() && dtEnd.Before(start) {
			continue
		}
		return true
	}
	return false
}
