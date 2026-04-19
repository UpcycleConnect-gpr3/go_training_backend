package training_content_models

import (
	"database/sql"
	"fmt"
	"go-training-backend/database"
	"go-training-backend/utils/log"
	"time"
)

const TABLE = "TRAINING_CONTENT"

type TrainingContent struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateTrainingContentDTO struct {
	Type    string
	Name    string
	Content string
}

type UpdateTrainingContentDTO struct {
	Type    string
	Name    string
	Content string
}

func GetAllTrainingContent(page, limit int) []TrainingContent {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Training.Query(
		"SELECT id, type, name, content, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []TrainingContent{}
	}
	defer rows.Close()

	items := []TrainingContent{}
	for rows.Next() {
		var tc TrainingContent
		if err := rows.Scan(&tc.Id, &tc.Type, &tc.Name, &tc.Content, &tc.CreatedAt, &tc.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		items = append(items, tc)
	}
	return items
}

func GetTrainingContentByID(id int) *TrainingContent {
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)
	tc := TrainingContent{}

	row := database.Training.QueryRow(
		"SELECT id, type, name, content, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)
	err := row.Scan(&tc.Id, &tc.Type, &tc.Name, &tc.Content, &tc.CreatedAt, &tc.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}
	return &tc
}

func CreateTrainingContent(dto CreateTrainingContentDTO) *TrainingContent {
	action := "INSERT INTO " + TABLE

	result, err := database.Training.Exec(
		"INSERT INTO "+TABLE+" (type, name, content) VALUES (?, ?, ?)",
		dto.Type, dto.Name, dto.Content,
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
	return GetTrainingContentByID(int(id))
}

func UpdateTrainingContent(id int, dto UpdateTrainingContentDTO) *TrainingContent {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Training.Exec(
		"UPDATE "+TABLE+" SET type = ?, name = ?, content = ? WHERE id = ?",
		dto.Type, dto.Name, dto.Content, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return GetTrainingContentByID(id)
}

func DeleteTrainingContent(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
