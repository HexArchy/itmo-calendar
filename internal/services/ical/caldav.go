package ical

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hexarchy/itmo-calendar/internal/entities"

	ics "github.com/arran4/golang-ical"
)

type Service struct{}

// New returns a new iCal service.
func New() *Service {
	return &Service{}
}

// Generate returns iCalendar data for the given schedule.
func (s *Service) Generate(_ context.Context, schedule []entities.DaySchedule) (*ics.Calendar, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetProductId("-//ITMO Calendar//EN")
	cal.SetVersion("2.0")
	cal.SetCalscale("GREGORIAN")
	cal.SetXWRCalName("ITMO Calendar")

	now := time.Now().UTC()

	for _, day := range schedule {
		for _, lesson := range day.Lessons {
			uid := fmt.Sprintf("%s-%s-%s@itmo-calendar",
				sanitizeUID(lesson.Subject),
				sanitizeUID(lesson.TeacherName),
				lesson.Start.UTC().Format("20060102T150405Z"))

			event := cal.AddEvent(uid)

			event.SetSummary(lesson.Subject)
			event.SetDtStampTime(now)
			event.SetCreatedTime(now)
			event.SetStartAt(lesson.Start.UTC())
			event.SetEndAt(lesson.End.UTC())

			descParts := []string{
				lesson.TeacherName,
				lesson.Type,
			}

			if lesson.Format != "" {
				descParts = append(descParts, fmt.Sprintf("Формат: %s", lesson.Format))
			}

			if lesson.Group != "" {
				descParts = append(descParts, fmt.Sprintf("Группа: %s", lesson.Group))
			}

			if lesson.Note != "" {
				descParts = append(descParts, fmt.Sprintf("Заметки: %s", lesson.Note))
			}

			if lesson.ZoomURL != "" {
				descParts = append(descParts, fmt.Sprintf("Zoom: %s", lesson.ZoomURL))
				event.AddProperty(ics.ComponentProperty(ics.PropertyUrl), lesson.ZoomURL)
			}

			event.SetDescription(strings.Join(descParts, "\n"))

			location := strings.TrimSpace(fmt.Sprintf("%s Аудитория: %s", lesson.Building, lesson.Room))
			if location != "" {
				event.SetLocation(location)
			}

			event.AddProperty(ics.ComponentProperty(ics.PropertyCategories), lesson.Type)
			event.SetStatus(ics.ObjectStatusConfirmed)
			event.SetTimeTransparency(ics.TransparencyOpaque)
		}
	}

	return cal, nil
}

// sanitizeUID removes special characters from strings to create valid UIDs.
func sanitizeUID(input string) string {
	replacer := strings.NewReplacer(
		" ", "-",
		":", "",
		";", "",
		",", "",
		"(", "",
		")", "",
		"<", "",
		">", "",
		"@", "-at-",
	)
	return replacer.Replace(input)
}
