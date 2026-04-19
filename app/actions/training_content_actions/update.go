package training_content_actions

import (
	"go-training-backend/app/models/training_content_models"
	"go-training-backend/utils/rules"
)

type UpdateTrainingContentDTO struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func UpdateTrainingContent(id int, dto UpdateTrainingContentDTO) ([]rules.ValidationError, *training_content_models.TrainingContent) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	content := training_content_models.UpdateTrainingContent(id, training_content_models.UpdateTrainingContentDTO{
		Type:    dto.Type,
		Name:    dto.Name,
		Content: dto.Content,
	})

	return nil, content
}
