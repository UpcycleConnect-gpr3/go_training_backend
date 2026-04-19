package content_schedule_actions

import (
	"go-training-backend/app/models/content_schedule_models"
	"go-training-backend/utils/rules"
)

type CreateContentScheduleDTO struct {
	Title             string `json:"title"`
	DayNumber         int    `json:"day_number"`
	Duration          string `json:"duration"`
	Description       string `json:"description"`
	Content           string `json:"content"`
	IsPractical       bool   `json:"is_practical"`
	ResourcesRequired string `json:"resources_required"`
	OrderPosition     int    `json:"order_position"`
	TrainingID        int    `json:"training_id"`
}

func CreateContentSchedule(dto CreateContentScheduleDTO) ([]rules.ValidationError, *content_schedule_models.ContentSchedule) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	rules.IntMinLength(dto.TrainingID, 1, "training_id", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	schedule := content_schedule_models.CreateContentSchedule(content_schedule_models.CreateContentScheduleDTO{
		Title:             dto.Title,
		DayNumber:         dto.DayNumber,
		Duration:          dto.Duration,
		Description:       dto.Description,
		Content:           dto.Content,
		IsPractical:       dto.IsPractical,
		ResourcesRequired: dto.ResourcesRequired,
		OrderPosition:     dto.OrderPosition,
		TrainingID:        dto.TrainingID,
	})

	return nil, schedule
}
