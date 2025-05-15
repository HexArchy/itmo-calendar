package entities

import (
	ics "github.com/arran4/golang-ical"
)

type CalDav struct {
	ISU  int64         `json:"isu"`
	ICal *ics.Calendar `json:"ical"`
}
