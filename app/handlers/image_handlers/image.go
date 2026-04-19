package image_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/image_actions"
	"go-training-backend/app/middleware/auth_middleware"
	"go-training-backend/app/models/image_models"
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

func GetImagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	page, limit := parsePage(r)
	images := image_models.GetAllImages(page, limit)
	response.NewSuccessData(w, images)
}

func GetImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	image := image_models.GetImageByID(id)
	if image == nil {
		response.NewErrorMessage(w, response.ErrImageNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, image)
}

func CreateImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userID := auth_middleware.GetUserId(r.Context())

	var dto image_actions.CreateImageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, image := image_actions.CreateImage(dto, userID)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, image)
}

func UpdateImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if image_models.GetImageByID(id) == nil {
		response.NewErrorMessage(w, response.ErrImageNotFound, http.StatusNotFound)
		return
	}

	var dto image_actions.UpdateImageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, image := image_actions.UpdateImage(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, image)
}

func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.NewErrorMessage(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if image_models.GetImageByID(id) == nil {
		response.NewErrorMessage(w, response.ErrImageNotFound, http.StatusNotFound)
		return
	}

	image_models.DeleteImage(id)
	response.NewSuccessMessage(w, "Image deleted")
}
