package training_handlers

import (
	"net/http"

	"go-training-backend/app/models/training_models"
	"go-training-backend/utils/jwt"
	"go-training-backend/utils/log"
	"go-training-backend/utils/request"
	"go-training-backend/utils/response"
)

// isResponsable : seul un responsable (administrator) peut valider/refuser.
func isResponsable(r *http.Request) bool {
	return jwt.RoleFromToken(r.Header.Get("Authorization")) == "administrator"
}

func setTrainingStatus(w http.ResponseWriter, r *http.Request, status string) {
	log.Api(r)

	if !isResponsable(r) {
		response.NewErrorMessage(w, response.ErrForbidden, http.StatusForbidden)
		return
	}

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var training training_models.Training
	if err := training.Get([]string{"id"}, "id = ?", id); err != nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	if err := training_models.SetStatus(id, status); err != nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusInternalServerError)
		return
	}

	response.NewSuccessData(w, map[string]any{"id": id, "status": status})
}

// ValidateTrainingHandler — POST /trainings/{id}/validate (responsable)
func ValidateTrainingHandler(w http.ResponseWriter, r *http.Request) {
	setTrainingStatus(w, r, "validated")
}

// RejectTrainingHandler — POST /trainings/{id}/reject (responsable)
func RejectTrainingHandler(w http.ResponseWriter, r *http.Request) {
	setTrainingStatus(w, r, "rejected")
}
