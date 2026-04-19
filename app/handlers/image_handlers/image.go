package image_handlers

import (
	"encoding/json"
	"go-training-backend/app/actions/image_actions"
	"go-training-backend/app/models/image_models"
	"go-training-backend/utils/auth"
	"go-training-backend/utils/log"
	"go-training-backend/utils/request"
	"go-training-backend/utils/response"
	"net/http"
)

func GetImagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	var image image_models.Image
	var images []image_models.Image

	columns := []string{"id", "path", "description", "created_by_user_id", "created_at", "updated_at"}

	err := image.All(columns, &images)
	if err != nil {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusInternalServerError)
		return
	}
	response.NewSuccessData(w, images)
}

func GetImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var image image_models.Image
	columns := []string{"id", "path", "description", "created_by_user_id", "created_at", "updated_at"}
	err := image.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrImageNotFound, http.StatusNotFound)
		return
	}

	response.NewSuccessData(w, image)
}

func CreateImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userID := auth.Auth(r).Id()

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

	response.NewSuccessData(w, map[string]int{"image_id": image.Id})
}

func UpdateImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var image image_models.Image
	columns := []string{"id"}
	err := image.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrImageNotFound, http.StatusNotFound)
		return
	}

	var dto image_actions.UpdateImageDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	validationErrors, updatedImage := image_actions.UpdateImage(id, dto)
	if len(validationErrors) > 0 {
		response.NewValidationError(w, response.ErrInvalidBody, validationErrors)
		return
	}

	response.NewSuccessData(w, map[string]int{"image_id": updatedImage.Id})
}

func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var image image_models.Image
	columns := []string{"id"}
	err := image.Get(columns, "id = ?", id)
	if err != nil {
		response.NewErrorMessage(w, response.ErrImageNotFound, http.StatusNotFound)
		return
	}

	image_models.DeleteImage(id)
	response.NewSuccessMessage(w, "Image deleted")
}
