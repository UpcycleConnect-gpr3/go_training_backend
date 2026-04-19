package database

import (
	"go-training-backend/config"
	"go-training-backend/database"
	"go-training-backend/internal"
	"go-training-backend/utils/log"

	"github.com/joho/godotenv"
)

func initialize() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Config Initialization
	config.InitDatabase()

	err = database.Training.Ping()

	if err != nil {
		log.Fatal(err)
	}

	internal.CreateTableMigrations(database.Training)

}

func Migrate() {

	initialize()

	internal.Migrate(database.Training)

}
