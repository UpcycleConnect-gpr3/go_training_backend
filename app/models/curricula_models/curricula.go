package curricula_models

import (
	"database/sql"
	"fmt"
	"go-training-backend/database"
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

func GetAllCurricula(page, limit int) []Curricula {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Training.Query(
		"SELECT id, path, name, description, created_by_user_id, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []Curricula{}
	}
	defer rows.Close()

	curricula := []Curricula{}
	for rows.Next() {
		var c Curricula
		if err := rows.Scan(&c.Id, &c.Path, &c.Name, &c.Description, &c.CreatedByUserID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		curricula = append(curricula, c)
	}
	return curricula
}

func GetCurriculaByID(id int) *Curricula {
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)
	c := Curricula{}

	row := database.Training.QueryRow(
		"SELECT id, path, name, description, created_by_user_id, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)
	err := row.Scan(&c.Id, &c.Path, &c.Name, &c.Description, &c.CreatedByUserID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}
	return &c
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
	return GetCurriculaByID(int(id))
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
	return GetCurriculaByID(id)
}

func DeleteCurricula(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
