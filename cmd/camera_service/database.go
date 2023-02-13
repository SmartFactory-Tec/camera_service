package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SmartFactory-Tec/camera_service/pkg/migrations"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"time"
)

type DbConfig struct {
	Hostname string
	Port     int
	Database string
	User     string
	Password string
}

func connectToDb(config DbConfig, logger *zap.SugaredLogger) *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.User, config.Hostname, config.Port, config.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("error parsing db connection string: %w", err)
	}

	testConnection(db, logger)

	logger.Infow("connected to database", "name", config.Database)

	return db
}

func testConnection(db *sql.DB, logger *zap.SugaredLogger) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {

		logger.Fatal(fmt.Errorf("could not connectToDb to database: %w", err))
	}
}

func updateDatabaseSchema(db *sql.DB, logger *zap.SugaredLogger) {
	goose.SetBaseFS(migrations.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatal(fmt.Errorf("could not initialize goose loader: %w", err))
	}

	if err := goose.Up(db, "."); err != nil {
		logger.Fatal(fmt.Errorf("failed to update database schema: %w", err))
	}

	logger.Infow("updated database schema")
}
