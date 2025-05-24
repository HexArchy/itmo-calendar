package itmoschedule

import (
	"time"

	"github.com/pkg/errors"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

type scheduleResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    []scheduleDayDTO `json:"data"`
}

type scheduleDayDTO struct {
	Date    string      `json:"date"`
	Lessons []lessonDTO `json:"lessons"`
}

type lessonDTO struct {
	Subject     string `json:"subject"`
	Type        string `json:"type"`
	TimeStart   string `json:"time_start"`
	TimeEnd     string `json:"time_end"`
	TeacherName string `json:"teacher_name"`
	Room        string `json:"room"`
	Note        string `json:"note"`
	Building    string `json:"building"`
	Format      string `json:"format"`
	Group       string `json:"group"`
	ZoomURL     string `json:"zoom_url"`
}

// transformDay converts a single day DTO to domain entity.
func (c *Client) transformDay(day scheduleDayDTO) (entities.DaySchedule, error) {
	date, err := time.Parse("2006-01-02", day.Date)
	if err != nil {
		return entities.DaySchedule{}, errors.Wrapf(err, "parse date %q", day.Date)
	}

	lessons := make([]entities.Lesson, 0, len(day.Lessons))
	for _, lesson := range day.Lessons {
		transformedLesson, err := c.transformLesson(day.Date, lesson)
		if err != nil {
			return entities.DaySchedule{}, errors.Wrap(err, "transform lesson")
		}

		lessons = append(lessons, transformedLesson)
	}

	return entities.DaySchedule{
		Date:    date,
		Lessons: lessons,
	}, nil
}

// transformLesson converts a lesson DTO to domain entity.
func (c *Client) transformLesson(dateStr string, lesson lessonDTO) (entities.Lesson, error) {
	startTime, err := parseDateTime(dateStr, lesson.TimeStart)
	if err != nil {
		return entities.Lesson{}, errors.Wrapf(err, "parse start time %q %q", dateStr, lesson.TimeStart)
	}

	endTime, err := parseDateTime(dateStr, lesson.TimeEnd)
	if err != nil {
		return entities.Lesson{}, errors.Wrapf(err, "parse end time %q %q", dateStr, lesson.TimeEnd)
	}

	// Handle potential null values.
	var note, zoomURL string
	if lesson.Note != "" {
		note = lesson.Note
	}
	if lesson.ZoomURL != "" {
		zoomURL = lesson.ZoomURL
	}

	return entities.Lesson{
		Subject:     lesson.Subject,
		Type:        lesson.Type,
		TeacherName: lesson.TeacherName,
		Room:        lesson.Room,
		Note:        note,
		Building:    lesson.Building,
		Format:      lesson.Format,
		Group:       lesson.Group,
		ZoomURL:     zoomURL,
		Start:       startTime,
		End:         endTime,
	}, nil
}

// parseDateTime combines date string (YYYY-MM-DD) and time string (HH:MM) into time.Time in Moscow time zone.
func parseDateTime(dateStr, timeStr string) (time.Time, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, errors.Wrap(err, "load Europe/Moscow location")
	}
	date, err := time.ParseInLocation("2006-01-02 15:04", dateStr+" "+timeStr, loc)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "parse date %q and time %q", dateStr, timeStr)
	}

	return date, nil
}
