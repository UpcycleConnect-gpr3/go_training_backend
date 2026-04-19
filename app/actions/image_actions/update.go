package image_actions

import (
	"go-training-backend/app/models/image_models"
	"go-training-backend/utils/rules"
)

type UpdateImageDTO struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

func UpdateImage(id int, dto UpdateImageDTO) ([]rules.ValidationError, *image_models.Image) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Path, 1, "path", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	image := image_models.UpdateImage(id, image_models.UpdateImageDTO{
		Path:        dto.Path,
		Description: dto.Description,
	})

	return nil, image
}
