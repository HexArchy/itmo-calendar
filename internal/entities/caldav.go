package entities

import (
	"errors"

	ics "github.com/arran4/golang-ical"
)

// ErrNotFound is returned when the requested entity does not exist.
var ErrNotFound = errors.New("not found")

type CalDav struct {
	ISU  int64         `json:"isu"`
	ICal *ics.Calendar `json:"ical"`
}
