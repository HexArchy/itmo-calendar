package ical

import (
	"context"
	"fmt"
	"regexp"
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

// Parse converts iCalendar data into DaySchedule entities.
func (s *Service) Parse(_ context.Context, cal *ics.Calendar) ([]entities.DaySchedule, error) {
	scheduleMap := make(map[string]*entities.DaySchedule)

	for _, event := range cal.Events() {
		lesson := entities.Lesson{}

		// Extract basic event information
		if summary := event.GetProperty(ics.ComponentPropertySummary); summary != nil {
			lesson.Subject = summary.Value
		}

		if dtstart := event.GetProperty(ics.ComponentPropertyDtStart); dtstart != nil {
			if startTime, err := time.Parse("20060102T150405Z", dtstart.Value); err == nil {
				lesson.Start = startTime
			}
		}

		if dtend := event.GetProperty(ics.ComponentPropertyDtEnd); dtend != nil {
			if endTime, err := time.Parse("20060102T150405Z", dtend.Value); err == nil {
				lesson.End = endTime
			}
		}

		// Parse description to extract structured data
		if desc := event.GetProperty(ics.ComponentPropertyDescription); desc != nil {
			s.parseDescription(desc.Value, &lesson)
		}

		// Extract location information
		if location := event.GetProperty(ics.ComponentPropertyLocation); location != nil {
			s.parseLocation(location.Value, &lesson)
		}

		// Extract type from categories
		if categories := event.GetProperty(ics.ComponentPropertyCategories); categories != nil {
			lesson.Type = categories.Value
		}

		// Extract Zoom URL from URL property
		if url := event.GetProperty(ics.ComponentPropertyUrl); url != nil {
			lesson.ZoomURL = url.Value
		}

		// Group lessons by date
		dateKey := lesson.Start.Format("2006-01-02")
		if _, exists := scheduleMap[dateKey]; !exists {
			date, _ := time.Parse("2006-01-02", dateKey)
			scheduleMap[dateKey] = &entities.DaySchedule{
				Date:    date,
				Lessons: []entities.Lesson{},
			}
		}

		scheduleMap[dateKey].Lessons = append(scheduleMap[dateKey].Lessons, lesson)
	}

	// Convert map to slice and sort by date
	var schedule []entities.DaySchedule
	for _, daySchedule := range scheduleMap {
		schedule = append(schedule, *daySchedule)
	}

	// Sort schedule by date (simple bubble sort for small datasets)
	for i := 0; i < len(schedule)-1; i++ {
		for j := 0; j < len(schedule)-i-1; j++ {
			if schedule[j].Date.After(schedule[j+1].Date) {
				schedule[j], schedule[j+1] = schedule[j+1], schedule[j]
			}
		}
	}

	return schedule, nil
}

// parseDescription extracts structured information from event description.
func (s *Service) parseDescription(description string, lesson *entities.Lesson) {
	lines := strings.Split(description, "\n")

	if len(lines) > 0 {
		lesson.TeacherName = strings.TrimSpace(lines[0])
	}

	if len(lines) > 1 {
		lesson.Type = strings.TrimSpace(lines[1])
	}

	// Parse additional fields using regex patterns
	formatRegex := regexp.MustCompile(`Формат:\s*(.+)`)
	groupRegex := regexp.MustCompile(`Группа:\s*(.+)`)
	noteRegex := regexp.MustCompile(`Заметки:\s*(.+)`)
	zoomRegex := regexp.MustCompile(`Zoom:\s*(.+)`)

	for _, line := range lines[2:] {
		line = strings.TrimSpace(line)

		if match := formatRegex.FindStringSubmatch(line); len(match) > 1 {
			lesson.Format = strings.TrimSpace(match[1])
		} else if match := groupRegex.FindStringSubmatch(line); len(match) > 1 {
			lesson.Group = strings.TrimSpace(match[1])
		} else if match := noteRegex.FindStringSubmatch(line); len(match) > 1 {
			lesson.Note = strings.TrimSpace(match[1])
		} else if match := zoomRegex.FindStringSubmatch(line); len(match) > 1 {
			lesson.ZoomURL = strings.TrimSpace(match[1])
		}
	}
}

// parseLocation extracts building and room information from location string.
func (s *Service) parseLocation(location string, lesson *entities.Lesson) {
	// Expected format: "Building Аудитория: Room"
	audienceRegex := regexp.MustCompile(`(.+?)\s*Аудитория:\s*(.+)`)

	if match := audienceRegex.FindStringSubmatch(location); len(match) > 2 {
		lesson.Building = strings.TrimSpace(match[1])
		lesson.Room = strings.TrimSpace(match[2])
	} else {
		// Fallback: treat entire location as building if pattern doesn't match
		lesson.Building = strings.TrimSpace(location)
	}
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
