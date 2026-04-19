package content_schedule_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/content_schedule_actions"
	"go-training-backend/app/models/content_schedule_models"
	"go-training-backend/utils/log"
	"go-training-backend/utils/request"
	"go-training-backend/utils/response"
	"net/http"
)

func GetContentSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var schedule content_schedule_models.ContentSchedule
	var schedules []content_schedule_models.ContentSchedule

	columns := []string{"id", "title", "day_number", "duration", "description", "content", "is_practical", "resources_required", "order_position", "training_id", "created_at", "updated_at"}

	err := schedule.All(columns, &schedules)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, schedules)
}

func GetContentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var schedule content_schedule_models.ContentSchedule
	columns := []string{"id", "title", "day_number", "duration", "description", "content", "is_practical", "resources_required", "order_position", "training_id", "created_at", "updated_at"}
	err := schedule.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrScheduleNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, schedule)
}

func CreateContentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto content_schedule_actions.CreateContentScheduleDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, schedule := content_schedule_actions.CreateContentSchedule(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, schedule)
}

func UpdateContentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var schedule content_schedule_models.ContentSchedule
	columns := []string{"id"}
	err := schedule.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrScheduleNotFound, http.StatusNotFound)
		return
	}

	var dto content_schedule_actions.UpdateContentScheduleDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, updatedSchedule := content_schedule_actions.UpdateContentSchedule(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, updatedSchedule)
}

func DeleteContentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var schedule content_schedule_models.ContentSchedule
	columns := []string{"id"}
	err := schedule.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrScheduleNotFound, http.StatusNotFound)
		return
	}

	content_schedule_models.DeleteContentSchedule(id)
	response.NewSuccessMessage(w, "Schedule deleted")
}
