package curricula_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/curricula_actions"
	"go-training-backend/app/models/curricula_models"
	"go-training-backend/utils/auth"
	"go-training-backend/utils/log"
	"go-training-backend/utils/request"
	"go-training-backend/utils/response"
	"net/http"
)

func GetCurriculaHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var curricula curricula_models.Curricula
	var curriculaList []curricula_models.Curricula

	columns := []string{"id", "path", "name", "description", "created_by_user_id", "created_at", "updated_at"}

	err := curricula.All(columns, &curriculaList)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, curriculaList)
}

func GetCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var curriculum curricula_models.Curricula
	columns := []string{"id", "path", "name", "description", "created_by_user_id", "created_at", "updated_at"}
	err := curriculum.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrCurriculaNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, curriculum)
}

func CreateCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userID := auth.Auth(r).Id()

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

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var curricula curricula_models.Curricula
	columns := []string{"id"}
	err := curricula.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrCurriculaNotFound, http.StatusNotFound)
		return
	}

	var dto curricula_actions.UpdateCurriculaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, updatedCurriculum := curricula_actions.UpdateCurricula(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, updatedCurriculum)
}

func DeleteCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var curricula curricula_models.Curricula
	columns := []string{"id"}
	err := curricula.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrCurriculaNotFound, http.StatusNotFound)
		return
	}

	curricula_models.DeleteCurricula(id)
	response.NewSuccessMessage(w, "Curricula deleted")
}
