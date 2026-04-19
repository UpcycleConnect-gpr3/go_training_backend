package category_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/category_actions"
	"go-training-backend/app/models/category_models"
	"go-training-backend/utils/auth"
	"go-training-backend/utils/log"
	"go-training-backend/utils/request"
	"go-training-backend/utils/response"
	"net/http"
)

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var category category_models.Category
	var categories []category_models.Category

	columns := []string{"id", "name", "slug", "description", "created_by_user_id", "created_at", "updated_at"}

	err := category.All(columns, &categories)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, categories)
}

func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var category category_models.Category
	columns := []string{"id", "name", "slug", "description", "created_by_user_id", "created_at", "updated_at"}
	err := category.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrCategoryNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, category)
}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userID := auth.Auth(r).Id()

	var dto category_actions.CreateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, category := category_actions.CreateCategory(dto, userID)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, category)
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var category category_models.Category
	columns := []string{"id"}
	err := category.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrCategoryNotFound, http.StatusNotFound)
		return
	}

	var dto category_actions.UpdateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, updatedCategory := category_actions.UpdateCategory(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, updatedCategory)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var category category_models.Category
	columns := []string{"id"}
	err := category.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrCategoryNotFound, http.StatusNotFound)
		return
	}

	category_models.DeleteCategory(id)
	response.NewSuccessMessage(w, "Category deleted")
}
