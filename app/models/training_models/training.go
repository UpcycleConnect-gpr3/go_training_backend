package training_models

import (
	"fmt"
	"go-training-backend/database"
	"go-training-backend/utils/db"
	"go-training-backend/utils/log"
)

const TABLE = "TRAININGS"

type Training struct {
	Id                          int     `json:"id"`
	Type                        string  `json:"type"`
	Name                        string  `json:"name"`
	ModeOfDelivery              string  `db:"mode_of_delivery" json:"mode_of_delivery"`
	Duration                    string  `json:"duration"`
	TargetAudience              string  `db:"target_audience" json:"target_audience"`
	MinimumNumberOfParticipants int     `db:"minimum_number_of_participants" json:"minimum_number_of_participants"`
	MaximumNumberOfParticipants int     `db:"maximum_number_of_participants" json:"maximum_number_of_participants"`
	Location                    string  `json:"location"`
	TrainerProfile              string  `db:"trainer_profile" json:"trainer_profile"`
	Price                       float64 `json:"price"`
	Status                      string  `json:"status"`
	CreatedAt                   string  `db:"created_at" json:"created_at"`
	UpdatedAt                   string  `db:"updated_at" json:"updated_at"`
}

type CreateTrainingDTO struct {
	Type                        string
	Name                        string
	ModeOfDelivery              string
	Duration                    string
	TargetAudience              string
	MinimumNumberOfParticipants int
	MaximumNumberOfParticipants int
	Location                    string
	TrainerProfile              string
	Price                       float64
}

type UpdateTrainingDTO struct {
	Type                        string
	Name                        string
	ModeOfDelivery              string
	Duration                    string
	TargetAudience              string
	MinimumNumberOfParticipants int
	MaximumNumberOfParticipants int
	Location                    string
	TrainerProfile              string
	Price                       float64
}

type CurriculaSummary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type TrainingContentSummary struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (t *Training) Get(columns []string, by string, value any) error {
	return db.GetQuery[Training](database.Training, TABLE, columns, by, value, t)
}

func (t *Training) All(columns []string, dest *[]Training) error {
	return db.AllQuery[Training](database.Training, TABLE, columns, dest)
}

// SetStatus met a jour le statut de validation d'une formation
// (pending / validated / rejected).
func SetStatus(id int, status string) error {
	_, err := database.Training.Exec(
		"UPDATE "+TABLE+" SET status = ?, updated_at = NOW() WHERE id = ?",
		status, id,
	)
	if err != nil {
		log.Database("SET TRAINING STATUS", err)
	}
	return err
}

func CreateTraining(dto CreateTrainingDTO) *Training {
	action := "INSERT INTO " + TABLE

	result, err := database.Training.Exec(
		"INSERT INTO "+TABLE+" (type, name, mode_of_delivery, duration, target_audience, minimum_number_of_participants, maximum_number_of_participants, location, trainer_profile, price) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		dto.Type, dto.Name, dto.ModeOfDelivery, dto.Duration, dto.TargetAudience,
		dto.MinimumNumberOfParticipants, dto.MaximumNumberOfParticipants,
		dto.Location, dto.TrainerProfile, dto.Price,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Training{Id: int(id)}
}

func UpdateTraining(id int, dto UpdateTrainingDTO) *Training {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Training.Exec(
		"UPDATE "+TABLE+" SET type = ?, name = ?, mode_of_delivery = ?, duration = ?, target_audience = ?, minimum_number_of_participants = ?, maximum_number_of_participants = ?, location = ?, trainer_profile = ?, price = ? WHERE id = ?",
		dto.Type, dto.Name, dto.ModeOfDelivery, dto.Duration, dto.TargetAudience,
		dto.MinimumNumberOfParticipants, dto.MaximumNumberOfParticipants,
		dto.Location, dto.TrainerProfile, dto.Price, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Training{Id: id}
}

func DeleteTraining(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}

func GetTrainingCurricula(trainingID int) []CurriculaSummary {
	action := fmt.Sprintf("SELECT CURRICULA JOIN TRAINING_CURRICULUM WHERE training_id : %d", trainingID)

	rows, err := database.Training.Query(
		"SELECT c.id, c.name, c.path FROM CURRICULA c JOIN TRAINING_CURRICULUM tc ON c.id = tc.curriculum_id WHERE tc.training_id = ?",
		trainingID,
	)
	if err != nil {
		log.Database(action, err)
		return []CurriculaSummary{}
	}
	defer rows.Close()

	curricula := []CurriculaSummary{}
	for rows.Next() {
		var c CurriculaSummary
		if err := rows.Scan(&c.Id, &c.Name, &c.Path); err != nil {
			log.Database(action, err)
			continue
		}
		curricula = append(curricula, c)
	}
	return curricula
}

func LinkCurriculum(trainingID, curriculumID int) {
	action := fmt.Sprintf("INSERT INTO TRAINING_CURRICULUM training_id:%d curriculum_id:%d", trainingID, curriculumID)
	_, err := database.Training.Exec(
		"INSERT IGNORE INTO TRAINING_CURRICULUM (training_id, curriculum_id) VALUES (?, ?)",
		trainingID, curriculumID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func UnlinkCurriculum(trainingID, curriculumID int) {
	action := fmt.Sprintf("DELETE FROM TRAINING_CURRICULUM training_id:%d curriculum_id:%d", trainingID, curriculumID)
	_, err := database.Training.Exec(
		"DELETE FROM TRAINING_CURRICULUM WHERE training_id = ? AND curriculum_id = ?",
		trainingID, curriculumID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func GetTrainingContent(trainingID int) []TrainingContentSummary {
	action := fmt.Sprintf("SELECT TRAINING_CONTENT JOIN TRAINING_TRAINING_CONTENT WHERE training_id : %d", trainingID)

	rows, err := database.Training.Query(
		"SELECT tc.id, tc.name, tc.type FROM TRAINING_CONTENT tc JOIN TRAINING_TRAINING_CONTENT ttc ON tc.id = ttc.training_content_id WHERE ttc.training_id = ?",
		trainingID,
	)
	if err != nil {
		log.Database(action, err)
		return []TrainingContentSummary{}
	}
	defer rows.Close()

	content := []TrainingContentSummary{}
	for rows.Next() {
		var c TrainingContentSummary
		if err := rows.Scan(&c.Id, &c.Name, &c.Type); err != nil {
			log.Database(action, err)
			continue
		}
		content = append(content, c)
	}
	return content
}

func LinkTrainingContent(trainingID, contentID int) {
	action := fmt.Sprintf("INSERT INTO TRAINING_TRAINING_CONTENT training_id:%d content_id:%d", trainingID, contentID)
	_, err := database.Training.Exec(
		"INSERT IGNORE INTO TRAINING_TRAINING_CONTENT (training_id, training_content_id) VALUES (?, ?)",
		trainingID, contentID,
	)
	if err != nil {
		log.Database(action, err)
	}
}

func UnlinkTrainingContent(trainingID, contentID int) {
	action := fmt.Sprintf("DELETE FROM TRAINING_TRAINING_CONTENT training_id:%d content_id:%d", trainingID, contentID)
	_, err := database.Training.Exec(
		"DELETE FROM TRAINING_TRAINING_CONTENT WHERE training_id = ? AND training_content_id = ?",
		trainingID, contentID,
	)
	if err != nil {
		log.Database(action, err)
	}
}
