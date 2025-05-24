// This file is safe to edit. Once it exists it will not be overwritten

package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/models"
	apiSchedule "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/restapi/operations/schedule"
)

func (h *Handler) GetScheduleHandler(params apiSchedule.GetScheduleParams) middleware.Responder {
	schedule, err := h.usecases.GetSchedule.Execute(params.HTTPRequest.Context(), params.Isu)
	if err != nil {
		return apiSchedule.NewGetScheduleInternalServerError().WithPayload(&models.Error{
			Error:   "InternalServerError",
			Message: err.Error(),
		})
	}
	if schedule == nil {
		return apiSchedule.NewGetScheduleNotFound().WithPayload(&models.Error{
			Error:   "NotFound",
			Message: "schedule not found",
		})
	}

	scheduleDTO := make([]*models.ScheduleItem, 0, len(schedule))

	for _, daySchedule := range schedule {
		date := strfmt.Date(daySchedule.Date)
		lessonsDTO := make([]*models.ScheduleItemLessonsItems0, 0, len(daySchedule.Lessons))
		for _, lesson := range daySchedule.Lessons {
			lessonDTO := &models.ScheduleItemLessonsItems0{
				Subject:     &lesson.Subject,
				Type:        &lesson.Type,
				TeacherName: &lesson.TeacherName,
				Room:        &lesson.Room,
				Note:        lesson.Note,
				Building:    &lesson.Building,
				Format:      &lesson.Format,
				Group:       &lesson.Group,
				ZoomURL:     lesson.ZoomURL,
				TimeStart:   (*strfmt.DateTime)(&lesson.Start),
				TimeEnd:     (*strfmt.DateTime)(&lesson.End),
			}
			lessonsDTO = append(lessonsDTO, lessonDTO)
		}

		scheduleItemDTO := &models.ScheduleItem{
			Date:    &date,
			Lessons: lessonsDTO,
		}

		scheduleDTO = append(scheduleDTO, scheduleItemDTO)
	}

	return apiSchedule.NewGetScheduleOK().WithPayload(scheduleDTO)
}
