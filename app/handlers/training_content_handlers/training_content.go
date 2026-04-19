package training_content_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/training_content_actions"
	"go-training-backend/app/models/training_content_models"
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

func GetTrainingContentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	content := training_content_models.GetAllTrainingContent(page, limit)
	response.NewSuccessData(w, content)
}

func GetTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	content := training_content_models.GetTrainingContentByID(id)
	if content == nil {
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

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if training_content_models.GetTrainingContentByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTrainingContentNotFound, http.StatusNotFound)
		return
	}

	var dto training_content_actions.UpdateTrainingContentDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, content := training_content_actions.UpdateTrainingContent(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, content)
}

func DeleteTrainingContentHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if training_content_models.GetTrainingContentByID(id) == nil {
		response.NewErrorMessage(w, response.ErrTrainingContentNotFound, http.StatusNotFound)
		return
	}

	training_content_models.DeleteTrainingContent(id)
	response.NewSuccessMessage(w, "Training content deleted")
}
