package entities

import "time"

// DaySchedule represents a day's schedule.
type DaySchedule struct {
	Date    time.Time `json:"date"`
	Lessons []Lesson  `json:"lessons"`
}

// Lesson represents a single lesson in the schedule.
type Lesson struct {
	Subject     string    `json:"subject"`
	Type        string    `json:"type"`
	TeacherName string    `json:"teacher_name"`
	Room        string    `json:"room"`
	Note        string    `json:"note,omitempty"`
	Building    string    `json:"building"`
	Format      string    `json:"format"`
	Group       string    `json:"group"`
	ZoomURL     string    `json:"zoom_url,omitempty"`
	Start       time.Time `json:"time_start"`
	End         time.Time `json:"time_end"`
}
