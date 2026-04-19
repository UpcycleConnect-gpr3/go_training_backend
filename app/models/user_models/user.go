package user_models

import (
	"go-training-backend/database"
	"go-training-backend/utils/db"

	"github.com/google/uuid"
)

const (
	TABLE = "USERS"
)

type User struct {
	Id        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Firstname string    `db:"firstname" json:"firstname"`
	Lastname  string    `db:"lastname" json:"lastname"`
	Email     string    `db:"email" json:"email"`
	CreatedAt string    `db:"created_at" json:"created_at"`
	UpdatedAt string    `db:"updated_at" json:"updated_at"`
}

func (user *User) Get(columns []string, by string, value any) error {
	return db.GetQuery[User](database.Training, TABLE, columns, by, value, user)
}

func (user *User) All(columns []string, dest *[]User) error {
	return db.AllQuery[User](database.Training, TABLE, columns, dest)
}
