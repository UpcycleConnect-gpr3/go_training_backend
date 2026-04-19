package curricula_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/curricula_actions"
	"go-training-backend/app/middleware/auth_middleware"
	"go-training-backend/app/models/curricula_models"
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

func GetCurriculaHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	curricula := curricula_models.GetAllCurricula(page, limit)
	response.NewSuccessData(w, curricula)
}

func GetCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	curricula := curricula_models.GetCurriculaByID(id)
	if curricula == nil {
		response.NewErrorMessage(w, response.ErrCurriculaNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, curricula)
}

func CreateCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userID := auth_middleware.GetUserId(r.Context())

	var dto curricula_actions.CreateCurriculaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, curricula := curricula_actions.CreateCurricula(dto, userID)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, curricula)
}

func UpdateCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if curricula_models.GetCurriculaByID(id) == nil {
		response.NewErrorMessage(w, response.ErrCurriculaNotFound, http.StatusNotFound)
		return
	}

	var dto curricula_actions.UpdateCurriculaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, curricula := curricula_actions.UpdateCurricula(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, curricula)
}

func DeleteCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if curricula_models.GetCurriculaByID(id) == nil {
		response.NewErrorMessage(w, response.ErrCurriculaNotFound, http.StatusNotFound)
		return
	}

	curricula_models.DeleteCurricula(id)
	response.NewSuccessMessage(w, "Curricula deleted")
}
