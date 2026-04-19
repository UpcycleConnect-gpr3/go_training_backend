package training_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/training_actions"
	"go-training-backend/app/models/content_schedule_models"
	"go-training-backend/app/models/training_models"
	"go-training-backend/utils/log"
	"go-training-backend/utils/request"
	"go-training-backend/utils/response"
	"net/http"
)

func GetTrainingsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var training training_models.Training
	var trainings []training_models.Training

	columns := []string{"id", "type", "name", "mode_of_delivery", "duration", "target_audience", "minimum_number_of_participants", "maximum_number_of_participants", "location", "trainer_profile", "created_at", "updated_at"}

	err := training.All(columns, &trainings)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, trainings)
}

func GetTrainingHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var training training_models.Training
	columns := []string{"id", "type", "name", "mode_of_delivery", "duration", "target_audience", "minimum_number_of_participants", "maximum_number_of_participants", "location", "trainer_profile", "created_at", "updated_at"}
	err := training.Get(columns, "id = ?", id)
	if err != nil {
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

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var training training_models.Training
	columns := []string{"id"}
	err := training.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	var dto training_actions.UpdateTrainingDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, updatedTraining := training_actions.UpdateTraining(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, updatedTraining)
}

func DeleteTrainingHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var training training_models.Training
	columns := []string{"id"}
	err := training.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	training_models.DeleteTraining(id)
	response.NewSuccessMessage(w, "Training deleted")
}

func GetTrainingCurriculaHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var training training_models.Training
	columns := []string{"id"}
	err := training.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	curricula := training_models.GetTrainingCurricula(id)
	response.NewSuccessData(w, curricula)
}

func LinkTrainingCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var training training_models.Training
	columns := []string{"id"}
	err := training.Get(columns, "id = ?", id)
	if err != nil {
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

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	curriculumID := request.Request(r, "curriculum_id").ConvertToInt(w)
	if curriculumID == -1 {
		return
	}

	training_models.UnlinkCurriculum(id, curriculumID)
	response.NewSuccessMessage(w, "Curriculum unlinked")
}

func GetTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var training training_models.Training
	columns := []string{"id"}
	err := training.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	content := training_models.GetTrainingContent(id)
	response.NewSuccessData(w, content)
}

func LinkTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var training training_models.Training
	columns := []string{"id"}
	err := training.Get(columns, "id = ?", id)
	if err != nil {
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

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	contentID := request.Request(r, "content_id").ConvertToInt(w)
	if contentID == -1 {
		return
	}

	training_models.UnlinkTrainingContent(id, contentID)
	response.NewSuccessMessage(w, "Content unlinked")
}

func GetTrainingSchedulesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var training training_models.Training
	columns := []string{"id"}
	err := training.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	schedules := content_schedule_models.GetSchedulesByTrainingID(id)
	response.NewSuccessData(w, schedules)
}
