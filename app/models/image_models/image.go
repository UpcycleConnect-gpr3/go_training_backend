package image_models

import (
	"fmt"
	"go-training-backend/database"
	"go-training-backend/utils/db"
	"go-training-backend/utils/log"
	"time"
)

const TABLE = "IMAGES"

type Image struct {
	Id              int       `json:"id"`
	Path            string    `json:"path"`
	Description     string    `json:"description"`
	CreatedByUserID string    `db:"created_by_user_id" json:"created_by_user_id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

func (img *Image) Get(columns []string, by string, value any) error {
	return db.GetQuery[Image](database.Training, TABLE, columns, by, value, img)
}

func (img *Image) All(columns []string, dest *[]Image) error {
	return db.AllQuery[Image](database.Training, TABLE, columns, dest)
}

type CreateImageDTO struct {
	Path            string
	Description     string
	CreatedByUserID string
}

type UpdateImageDTO struct {
	Path        string
	Description string
}

func CreateImage(dto CreateImageDTO) *Image {
	action := "INSERT INTO " + TABLE

	result, err := database.Training.Exec(
		"INSERT INTO "+TABLE+" (path, description, created_by_user_id) VALUES (?, ?, ?)",
		dto.Path, dto.Description, dto.CreatedByUserID,
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
	return &Image{Id: int(id)}
}

func UpdateImage(id int, dto UpdateImageDTO) *Image {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Training.Exec(
		"UPDATE "+TABLE+" SET path = ?, description = ? WHERE id = ?",
		dto.Path, dto.Description, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return &Image{Id: id}
}

func DeleteImage(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
