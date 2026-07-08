package training_actions

import (
	"go-training-backend/app/models/training_models"
	"go-training-backend/utils/rules"
)

type CreateTrainingDTO struct {
	Type                        string  `json:"type"`
	Name                        string  `json:"name"`
	ModeOfDelivery              string  `json:"mode_of_delivery"`
	Duration                    string  `json:"duration"`
	TargetAudience              string  `json:"target_audience"`
	MinimumNumberOfParticipants int     `json:"minimum_number_of_participants"`
	MaximumNumberOfParticipants int     `json:"maximum_number_of_participants"`
	Location                    string  `json:"location"`
	TrainerProfile              string  `json:"trainer_profile"`
	Price                       float64 `json:"price"`
}

func CreateTraining(dto CreateTrainingDTO) ([]rules.ValidationError, *training_models.Training) {
	var errs []rules.ValidationError

	rules.StringMinLength(dto.Name, 1, "name", &errs)
	if len(errs) > 0 {
		return errs, nil
	}

	training := training_models.CreateTraining(training_models.CreateTrainingDTO{
		Type:                        dto.Type,
		Name:                        dto.Name,
		ModeOfDelivery:              dto.ModeOfDelivery,
		Duration:                    dto.Duration,
		TargetAudience:              dto.TargetAudience,
		MinimumNumberOfParticipants: dto.MinimumNumberOfParticipants,
		MaximumNumberOfParticipants: dto.MaximumNumberOfParticipants,
		Location:                    dto.Location,
		TrainerProfile:              dto.TrainerProfile,
		Price:                       dto.Price,
	})

	return nil, training
}
