package content_schedule_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/content_schedule_actions"
	"go-training-backend/app/models/content_schedule_models"
	"go-training-backend/utils/log"
	"go-training-backend/utils/response"
	"net/http"
	"strconv"
)

func parsePage(r *http.Request) (int, int) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 20
	}
	return page, limit
}

func GetContentSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	schedules := content_schedule_models.GetAllContentSchedules(page, limit)
	response.NewSuccessData(w, schedules)
}

func GetContentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	schedule := content_schedule_models.GetContentScheduleByID(id)
	if schedule == nil {
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

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if content_schedule_models.GetContentScheduleByID(id) == nil {
		response.NewErrorMessage(w, response.ErrScheduleNotFound, http.StatusNotFound)
		return
	}

	var dto content_schedule_actions.UpdateContentScheduleDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, schedule := content_schedule_actions.UpdateContentSchedule(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, schedule)
}

func DeleteContentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if content_schedule_models.GetContentScheduleByID(id) == nil {
		response.NewErrorMessage(w, response.ErrScheduleNotFound, http.StatusNotFound)
		return
	}

	content_schedule_models.DeleteContentSchedule(id)
	response.NewSuccessMessage(w, "Schedule deleted")
}
