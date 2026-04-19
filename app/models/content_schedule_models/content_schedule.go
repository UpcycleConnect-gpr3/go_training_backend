package content_schedule_models

import (
	"database/sql"
	"fmt"
	"go-training-backend/database"
	"go-training-backend/utils/log"
	"time"
)

const TABLE = "CONTENT_AND_SCHEDULES"

type ContentSchedule struct {
	Id                int       `json:"id"`
	Title             string    `json:"title"`
	DayNumber         int       `db:"day_number" json:"day_number"`
	Duration          string    `json:"duration"`
	Description       string    `json:"description"`
	Content           string    `json:"content"`
	IsPractical       bool      `db:"is_practical" json:"is_practical"`
	ResourcesRequired string    `db:"resources_required" json:"resources_required"`
	OrderPosition     int       `db:"order_position" json:"order_position"`
	TrainingID        int       `db:"training_id" json:"training_id"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type CreateContentScheduleDTO struct {
	Title             string
	DayNumber         int
	Duration          string
	Description       string
	Content           string
	IsPractical       bool
	ResourcesRequired string
	OrderPosition     int
	TrainingID        int
}

type UpdateContentScheduleDTO struct {
	Title             string
	DayNumber         int
	Duration          string
	Description       string
	Content           string
	IsPractical       bool
	ResourcesRequired string
	OrderPosition     int
}

func GetAllContentSchedules(page, limit int) []ContentSchedule {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Training.Query(
		"SELECT id, title, day_number, duration, description, content, is_practical, resources_required, order_position, training_id, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []ContentSchedule{}
	}
	defer rows.Close()

	schedules := []ContentSchedule{}
	for rows.Next() {
		var s ContentSchedule
		if err := rows.Scan(
			&s.Id, &s.Title, &s.DayNumber, &s.Duration, &s.Description,
			&s.Content, &s.IsPractical, &s.ResourcesRequired, &s.OrderPosition,
			&s.TrainingID, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			log.Database(action, err)
			continue
		}
		schedules = append(schedules, s)
	}
	return schedules
}

func GetContentScheduleByID(id int) *ContentSchedule {
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)
	s := ContentSchedule{}

	row := database.Training.QueryRow(
		"SELECT id, title, day_number, duration, description, content, is_practical, resources_required, order_position, training_id, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)
	err := row.Scan(
		&s.Id, &s.Title, &s.DayNumber, &s.Duration, &s.Description,
		&s.Content, &s.IsPractical, &s.ResourcesRequired, &s.OrderPosition,
		&s.TrainingID, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}
	return &s
}

func GetSchedulesByTrainingID(trainingID int) []ContentSchedule {
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE training_id : %d", trainingID)

	rows, err := database.Training.Query(
		"SELECT id, title, day_number, duration, description, content, is_practical, resources_required, order_position, training_id, created_at, updated_at FROM "+TABLE+" WHERE training_id = ? ORDER BY day_number, order_position",
		trainingID,
	)
	if err != nil {
		log.Database(action, err)
		return []ContentSchedule{}
	}
	defer rows.Close()

	schedules := []ContentSchedule{}
	for rows.Next() {
		var s ContentSchedule
		if err := rows.Scan(
			&s.Id, &s.Title, &s.DayNumber, &s.Duration, &s.Description,
			&s.Content, &s.IsPractical, &s.ResourcesRequired, &s.OrderPosition,
			&s.TrainingID, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			log.Database(action, err)
			continue
		}
		schedules = append(schedules, s)
	}
	return schedules
}

func CreateContentSchedule(dto CreateContentScheduleDTO) *ContentSchedule {
	action := "INSERT INTO " + TABLE

	result, err := database.Training.Exec(
		"INSERT INTO "+TABLE+" (title, day_number, duration, description, content, is_practical, resources_required, order_position, training_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		dto.Title, dto.DayNumber, dto.Duration, dto.Description, dto.Content,
		dto.IsPractical, dto.ResourcesRequired, dto.OrderPosition, dto.TrainingID,
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
	return GetContentScheduleByID(int(id))
}

func UpdateContentSchedule(id int, dto UpdateContentScheduleDTO) *ContentSchedule {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Training.Exec(
		"UPDATE "+TABLE+" SET title = ?, day_number = ?, duration = ?, description = ?, content = ?, is_practical = ?, resources_required = ?, order_position = ? WHERE id = ?",
		dto.Title, dto.DayNumber, dto.Duration, dto.Description, dto.Content,
		dto.IsPractical, dto.ResourcesRequired, dto.OrderPosition, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return GetContentScheduleByID(id)
}

func DeleteContentSchedule(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
