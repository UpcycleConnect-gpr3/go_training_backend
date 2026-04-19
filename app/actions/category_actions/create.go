package category_actions

import (
	"go-training-backend/app/models/category_models"
	"go-training-backend/utils/rules"
)

type CreateCategoryDTO struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

func CreateCategory(dto CreateCategoryDTO, userID string) ([]rules.ValidationError, *category_models.Category) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	rules.StringMinLength(dto.Slug, 1, "slug", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	category := category_models.CreateCategory(category_models.CreateCategoryDTO{
		Name:            dto.Name,
		Slug:            dto.Slug,
		Description:     dto.Description,
		CreatedByUserID: userID,
	})

	return nil, category
}
