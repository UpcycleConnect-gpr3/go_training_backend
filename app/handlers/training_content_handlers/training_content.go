package training_content_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/training_content_actions"
	"go-training-backend/app/models/training_content_models"
	"go-training-backend/utils/log"
	"go-training-backend/utils/request"
	"go-training-backend/utils/response"
	"net/http"
)

func GetTrainingContentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var content training_content_models.TrainingContent
	var contents []training_content_models.TrainingContent

	columns := []string{"id", "type", "name", "content", "created_at", "updated_at"}

	err := content.All(columns, &contents)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, contents)
}

func GetTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var content training_content_models.TrainingContent
	columns := []string{"id", "type", "name", "content", "created_at", "updated_at"}
	err := content.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrTrainingContentNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, content)
}

func CreateTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var dto training_content_actions.CreateTrainingContentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, content := training_content_actions.CreateTrainingContent(dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, content)
}

func UpdateTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var content training_content_models.TrainingContent
	columns := []string{"id"}
	err := content.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrTrainingContentNotFound, http.StatusNotFound)
		return
	}

	var dto training_content_actions.UpdateTrainingContentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, updatedContent := training_content_actions.UpdateTrainingContent(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, updatedContent)
}

func DeleteTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var content training_content_models.TrainingContent
	columns := []string{"id"}
	err := content.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrTrainingContentNotFound, http.StatusNotFound)
		return
	}

	training_content_models.DeleteTrainingContent(id)
	response.NewSuccessMessage(w, "Training content deleted")
}
