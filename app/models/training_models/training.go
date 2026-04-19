package training_models

import (
	"database/sql"
	"fmt"
	"go-training-backend/database"
	"go-training-backend/utils/log"
	"time"
)

const TABLE = "TRAININGS"

type Training struct {
	Id                          int       `json:"id"`
	Type                        string    `json:"type"`
	Name                        string    `json:"name"`
	ModeOfDelivery              string    `db:"mode_of_delivery" json:"mode_of_delivery"`
	Duration                    string    `json:"duration"`
	TargetAudience              string    `db:"target_audience" json:"target_audience"`
	MinimumNumberOfParticipants int       `db:"minimum_number_of_participants" json:"minimum_number_of_participants"`
	MaximumNumberOfParticipants int       `db:"maximum_number_of_participants" json:"maximum_number_of_participants"`
	Location                    string    `json:"location"`
	TrainerProfile              string    `db:"trainer_profile" json:"trainer_profile"`
	CreatedAt                   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                   time.Time `db:"updated_at" json:"updated_at"`
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

func scanTraining(row *sql.Row) (*Training, error) {
	t := Training{}
	err := row.Scan(
		&t.Id, &t.Type, &t.Name, &t.ModeOfDelivery, &t.Duration,
		&t.TargetAudience, &t.MinimumNumberOfParticipants,
		&t.MaximumNumberOfParticipants, &t.Location, &t.TrainerProfile,
		&t.CreatedAt, &t.UpdatedAt,
	)
	return &t, err
}

func GetAllTrainings(page, limit int) []Training {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Training.Query(
		"SELECT id, type, name, mode_of_delivery, duration, target_audience, minimum_number_of_participants, maximum_number_of_participants, location, trainer_profile, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []Training{}
	}
	defer rows.Close()

	trainings := []Training{}
	for rows.Next() {
		var t Training
		if err := rows.Scan(
			&t.Id, &t.Type, &t.Name, &t.ModeOfDelivery, &t.Duration,
			&t.TargetAudience, &t.MinimumNumberOfParticipants,
			&t.MaximumNumberOfParticipants, &t.Location, &t.TrainerProfile,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			log.Database(action, err)
			continue
		}
		trainings = append(trainings, t)
	}
	return trainings
}

func GetTrainingByID(id int) *Training {
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)

	row := database.Training.QueryRow(
		"SELECT id, type, name, mode_of_delivery, duration, target_audience, minimum_number_of_participants, maximum_number_of_participants, location, trainer_profile, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)
	t, err := scanTraining(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}
	return t
}

func CreateTraining(dto CreateTrainingDTO) *Training {
	action := "INSERT INTO " + TABLE

	result, err := database.Training.Exec(
		"INSERT INTO "+TABLE+" (type, name, mode_of_delivery, duration, target_audience, minimum_number_of_participants, maximum_number_of_participants, location, trainer_profile) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		dto.Type, dto.Name, dto.ModeOfDelivery, dto.Duration, dto.TargetAudience,
		dto.MinimumNumberOfParticipants, dto.MaximumNumberOfParticipants,
		dto.Location, dto.TrainerProfile,
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
	return GetTrainingByID(int(id))
}

func UpdateTraining(id int, dto UpdateTrainingDTO) *Training {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Training.Exec(
		"UPDATE "+TABLE+" SET type = ?, name = ?, mode_of_delivery = ?, duration = ?, target_audience = ?, minimum_number_of_participants = ?, maximum_number_of_participants = ?, location = ?, trainer_profile = ? WHERE id = ?",
		dto.Type, dto.Name, dto.ModeOfDelivery, dto.Duration, dto.TargetAudience,
		dto.MinimumNumberOfParticipants, dto.MaximumNumberOfParticipants,
		dto.Location, dto.TrainerProfile, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return GetTrainingByID(id)
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
