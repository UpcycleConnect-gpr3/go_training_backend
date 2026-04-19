package image_models

import (
	"database/sql"
	"fmt"
	"go-training-backend/database"
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

type CreateImageDTO struct {
	Path            string
	Description     string
	CreatedByUserID string
}

type UpdateImageDTO struct {
	Path        string
	Description string
}

func GetAllImages(page, limit int) []Image {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Training.Query(
		"SELECT id, path, description, created_by_user_id, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []Image{}
	}
	defer rows.Close()

	images := []Image{}
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.Id, &img.Path, &img.Description, &img.CreatedByUserID, &img.CreatedAt, &img.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		images = append(images, img)
	}
	return images
}

func GetImageByID(id int) *Image {
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)
	img := Image{}

	row := database.Training.QueryRow(
		"SELECT id, path, description, created_by_user_id, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)
	err := row.Scan(&img.Id, &img.Path, &img.Description, &img.CreatedByUserID, &img.CreatedAt, &img.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}
	return &img
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
	return GetImageByID(int(id))
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
	return GetImageByID(id)
}

func DeleteImage(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
