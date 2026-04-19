package curricula_actions

import (
	"go-training-backend/app/models/curricula_models"
	"go-training-backend/utils/rules"
)

type CreateCurriculaDTO struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateCurricula(dto CreateCurriculaDTO, userID string) ([]rules.ValidationError, *curricula_models.Curricula) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	curricula := curricula_models.CreateCurricula(curricula_models.CreateCurriculaDTO{
		Path:            dto.Path,
		Name:            dto.Name,
		Description:     dto.Description,
		CreatedByUserID: userID,
	})

	return nil, curricula
}
