package training_content_models

import (
	"fmt"
	"go-training-backend/database"
	"go-training-backend/utils/db"
	"go-training-backend/utils/log"
)

const TABLE = "TRAINING_CONTENT"

type TrainingContent struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
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

func (tc *TrainingContent) Get(columns []string, by string, value any) error {
	return db.GetQuery[TrainingContent](database.Training, TABLE, columns, by, value, tc)
}

func (tc *TrainingContent) All(columns []string, dest *[]TrainingContent) error {
	return db.AllQuery[TrainingContent](database.Training, TABLE, columns, dest)
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
	return &TrainingContent{Id: int(id)}
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
	return &TrainingContent{Id: id}
}

func DeleteTrainingContent(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
