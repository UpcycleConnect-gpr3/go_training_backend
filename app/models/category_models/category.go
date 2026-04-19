package category_models

import (
	"database/sql"
	"fmt"
	"go-training-backend/database"
	"go-training-backend/utils/log"
	"time"
)

const TABLE = "CATEGORIES"

type Category struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     string    `json:"description"`
	CreatedByUserID string    `db:"created_by_user_id" json:"created_by_user_id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type CreateCategoryDTO struct {
	Name            string
	Slug            string
	Description     string
	CreatedByUserID string
}

type UpdateCategoryDTO struct {
	Name        string
	Slug        string
	Description string
}

func GetAllCategories(page, limit int) []Category {
	action := "SELECT " + TABLE + " (paginated)"
	offset := (page - 1) * limit

	rows, err := database.Training.Query(
		"SELECT id, name, slug, description, created_by_user_id, created_at, updated_at FROM "+TABLE+" LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Database(action, err)
		return []Category{}
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.Id, &c.Name, &c.Slug, &c.Description, &c.CreatedByUserID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Database(action, err)
			continue
		}
		categories = append(categories, c)
	}
	return categories
}

func GetCategoryByID(id int) *Category {
	action := fmt.Sprintf("SELECT "+TABLE+" WHERE id : %d", id)
	c := Category{}

	row := database.Training.QueryRow(
		"SELECT id, name, slug, description, created_by_user_id, created_at, updated_at FROM "+TABLE+" WHERE id = ?",
		id,
	)
	err := row.Scan(&c.Id, &c.Name, &c.Slug, &c.Description, &c.CreatedByUserID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Database(action, err)
		return nil
	}
	return &c
}

func CreateCategory(dto CreateCategoryDTO) *Category {
	action := "INSERT INTO " + TABLE

	result, err := database.Training.Exec(
		"INSERT INTO "+TABLE+" (name, slug, description, created_by_user_id) VALUES (?, ?, ?, ?)",
		dto.Name, dto.Slug, dto.Description, dto.CreatedByUserID,
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
	return GetCategoryByID(int(id))
}

func UpdateCategory(id int, dto UpdateCategoryDTO) *Category {
	action := fmt.Sprintf("UPDATE "+TABLE+" WHERE id : %d", id)

	_, err := database.Training.Exec(
		"UPDATE "+TABLE+" SET name = ?, slug = ?, description = ? WHERE id = ?",
		dto.Name, dto.Slug, dto.Description, id,
	)
	if err != nil {
		log.Database(action, err)
		return nil
	}
	return GetCategoryByID(id)
}

func DeleteCategory(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
