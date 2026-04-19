package content_schedule_actions

import (
	"go-training-backend/app/models/content_schedule_models"
	"go-training-backend/utils/rules"
)

type UpdateContentScheduleDTO struct {
	Title             string `json:"title"`
	DayNumber         int    `json:"day_number"`
	Duration          string `json:"duration"`
	Description       string `json:"description"`
	Content           string `json:"content"`
	IsPractical       bool   `json:"is_practical"`
	ResourcesRequired string `json:"resources_required"`
	OrderPosition     int    `json:"order_position"`
}

func UpdateContentSchedule(id int, dto UpdateContentScheduleDTO) ([]rules.ValidationError, *content_schedule_models.ContentSchedule) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Title, 1, "title", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	schedule := content_schedule_models.UpdateContentSchedule(id, content_schedule_models.UpdateContentScheduleDTO{
		Title:             dto.Title,
		DayNumber:         dto.DayNumber,
		Duration:          dto.Duration,
		Description:       dto.Description,
		Content:           dto.Content,
		IsPractical:       dto.IsPractical,
		ResourcesRequired: dto.ResourcesRequired,
		OrderPosition:     dto.OrderPosition,
	})

	return nil, schedule
}
