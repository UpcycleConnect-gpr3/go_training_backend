package image_actions

import (
	"go-training-backend/app/models/image_models"
	"go-training-backend/utils/rules"
)

type CreateImageDTO struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

func CreateImage(dto CreateImageDTO, userID string) ([]rules.ValidationError, *image_models.Image) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Path, 1, "path", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	image := image_models.CreateImage(image_models.CreateImageDTO{
		Path:            dto.Path,
		Description:     dto.Description,
		CreatedByUserID: userID,
	})

	return nil, image
}
