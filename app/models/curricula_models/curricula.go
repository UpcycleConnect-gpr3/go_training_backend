package curricula_models

import (
	"fmt"
	"go-training-backend/database"
	"go-training-backend/utils/db"
	"go-training-backend/utils/log"
	"time"
)

const TABLE = "CURRICULA"

type Curricula struct {
	Id              int       `json:"id"`
	Path            string    `json:"path"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CreatedByUserID string    `db:"created_by_user_id" json:"created_by_user_id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

func (c *Curricula) Get(columns []string, by string, value any) error {
	return db.GetQuery[Curricula](database.Training, TABLE, columns, by, value, c)
}

func (c *Curricula) All(columns []string, dest *[]Curricula) error {
	return db.AllQuery[Curricula](database.Training, TABLE, columns, dest)
}

type CreateCurriculaDTO struct {
	Path            string
	Name            string
	Description     string
	CreatedByUserID string
}

type UpdateCurriculaDTO struct {
	Path        string
	Name        string
	Description string
}

func CreateCurricula(dto CreateCurriculaDTO) *Curricula {
	action := "INSERT INTO " + TABLE

	result, err := database.Training.Exec(
		"INSERT INTO "+TABLE+" (path, name, description, created_by_user_id) VALUES (?, ?, ?, ?)",
		dto.Path, dto.Name, dto.Description, dto.CreatedByUserID,
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
	return &Curricula{Id: int(id)}
}

func UpdateCurricula(id int, dto UpdateCurriculaDTO) *Curricula {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Training.Exec(
		"UPDATE "+TABLE+" SET path = ?, name = ?, description = ? WHERE id = ?",
		dto.Path, dto.Name, dto.Description, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Curricula{Id: id}
}

func DeleteCurricula(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
