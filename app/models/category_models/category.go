package category_models

import (
	"fmt"
	"go-training-backend/database"
	"go-training-backend/utils/db"
	"go-training-backend/utils/log"
)

const TABLE = "CATEGORIES"

type Category struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Description     string `json:"description"`
	CreatedByUserID string `db:"created_by_user_id" json:"created_by_user_id"`
	CreatedAt       string `db:"created_at" json:"created_at"`
	UpdatedAt       string `db:"updated_at" json:"updated_at"`
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

func (category *Category) Get(columns []string, by string, value any) error {
	return db.GetQuery[Category](database.Training, TABLE, columns, by, value, category)
}

func (category *Category) All(columns []string, dest *[]Category) error {
	return db.AllQuery[Category](database.Training, TABLE, columns, dest)
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
	return &Category{Id: int(id)}
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
	return &Category{Id: int(id)}
}

func DeleteCategory(id int) {
	action := fmt.Sprintf("DELETE FROM "+TABLE+" WHERE id : %d", id)
	_, err := database.Training.Exec("DELETE FROM "+TABLE+" WHERE id = ?", id)
	if err != nil {
		log.Database(action, err)
	}
}
