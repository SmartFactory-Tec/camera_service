package main

import (
	_ "github.com/lib/pq"
)

func main() {
	logger := setupLogger()

	config := loadConfig(logger)

	dbConfig := config.Db

	db := connectToDb(dbConfig, logger)

	updateDatabaseSchema(db, logger)

}
