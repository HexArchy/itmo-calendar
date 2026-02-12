package entities

import (
	"errors"

	ics "github.com/arran4/golang-ical"
)

// ErrNotFound is returned when the requested entity does not exist.
var ErrNotFound = errors.New("not found")

// ErrAuthFailed is returned when authentication fails.
var ErrAuthFailed = errors.New("authentication failed")

// ErrUpstreamFail is returned when an upstream service is unavailable.
var ErrUpstreamFail = errors.New("upstream service unavailable")

type CalDav struct {
	ISU  int64         `json:"isu"`
	ICal *ics.Calendar `json:"ical"`
}
