package itmoschedule

import (
	"time"

	"github.com/pkg/errors"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

var moscowLocation = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic("failed to load Europe/Moscow location: " + err.Error())
	}
	return loc
}()

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
		transformedLesson, errTransform := c.transformLesson(date, lesson)
		if errTransform != nil {
			return entities.DaySchedule{}, errors.Wrap(errTransform, "transform lesson")
		}

		lessons = append(lessons, transformedLesson)
	}

	return entities.DaySchedule{
		Date:    date,
		Lessons: lessons,
	}, nil
}

// transformLesson converts a lesson DTO to domain entity.
func (c *Client) transformLesson(date time.Time, lesson lessonDTO) (entities.Lesson, error) {
	startTime, err := parseTimeOnDate(date, lesson.TimeStart)
	if err != nil {
		return entities.Lesson{}, errors.Wrapf(err, "parse start time %q", lesson.TimeStart)
	}

	endTime, err := parseTimeOnDate(date, lesson.TimeEnd)
	if err != nil {
		return entities.Lesson{}, errors.Wrapf(err, "parse end time %q", lesson.TimeEnd)
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

// parseTimeOnDate combines date and time string (HH:MM) into [time.Time] in Moscow time zone.
func parseTimeOnDate(date time.Time, timeStr string) (time.Time, error) {
	t, err := time.ParseInLocation("15:04", timeStr, moscowLocation)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "parse time %q", timeStr)
	}

	result := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		t.Hour(),
		t.Minute(),
		0,
		0,
		moscowLocation,
	)

	return result, nil
}
