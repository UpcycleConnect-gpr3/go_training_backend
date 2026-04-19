package training_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/training_actions"
	"go-training-backend/app/models/content_schedule_models"
	"go-training-backend/app/models/training_models"
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

func GetTrainingsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	trainings := training_models.GetAllTrainings(page, limit)
	response.NewSuccessData(w, trainings)
}

func GetTrainingHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	training := training_models.GetTrainingByID(id)
	if training == nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, training)
}

func CreateTrainingHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto training_actions.CreateTrainingDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, training := training_actions.CreateTraining(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, training)
}

func UpdateTrainingHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if training_models.GetTrainingByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	var dto training_actions.UpdateTrainingDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, training := training_actions.UpdateTraining(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, training)
}

func DeleteTrainingHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if training_models.GetTrainingByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	training_models.DeleteTraining(id)
	response.NewSuccessMessage(w, "Training deleted")
}

func GetTrainingCurriculaHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if training_models.GetTrainingByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	curricula := training_models.GetTrainingCurricula(id)
	response.NewSuccessData(w, curricula)
}

func LinkTrainingCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if training_models.GetTrainingByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	var body struct {
		CurriculumID int `json:"curriculum_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	training_models.LinkCurriculum(id, body.CurriculumID)
	response.NewSuccessMessage(w, "Curriculum linked")
}

func UnlinkTrainingCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	curriculumID, err := strconv.Atoi(r.PathValue("curriculum_id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid curriculum_id", http.StatusBadRequest)
		return
	}

	training_models.UnlinkCurriculum(id, curriculumID)
	response.NewSuccessMessage(w, "Curriculum unlinked")
}

func GetTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if training_models.GetTrainingByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	content := training_models.GetTrainingContent(id)
	response.NewSuccessData(w, content)
}

func LinkTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if training_models.GetTrainingByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	var body struct {
		ContentID int `json:"content_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	training_models.LinkTrainingContent(id, body.ContentID)
	response.NewSuccessMessage(w, "Content linked")
}

func UnlinkTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	contentID, err := strconv.Atoi(r.PathValue("content_id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid content_id", http.StatusBadRequest)
		return
	}

	training_models.UnlinkTrainingContent(id, contentID)
	response.NewSuccessMessage(w, "Content unlinked")
}

func GetTrainingSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if training_models.GetTrainingByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	schedules := content_schedule_models.GetSchedulesByTrainingID(id)
	response.NewSuccessData(w, schedules)
}
