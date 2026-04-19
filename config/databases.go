package config

import (
	"go-training-backend/database"
	"go-training-backend/internal"
	"os"
)

func InitDatabase() {

	database.Training = internal.NewDatabase(
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
}
