package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SmartFactory-Tec/camera_service/pkg/migrations"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"net/http"
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

		logger.Fatalf("could not connectToDb to database: %s", err)
	}
}

func updateDatabaseSchema(db *sql.DB, logger *zap.SugaredLogger) {
	goose.SetBaseFS(migrations.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatalf("could not initialize goose loader: %s", err)
	}

	if err := goose.Up(db, "."); err != nil {
		logger.Fatalf("failed to update database schema: %s", err)
	}

	logger.Infow("updated database schema")
}

func HandlePqError(w http.ResponseWriter, r *http.Request, err *pq.Error, logger *zap.SugaredLogger) {
	switch err.Code {
	case "23502":
		// not-null constraint violation
		logger.Errorf("not null constraint violation: %s", err.Message)
		http.Error(w, fmt.Sprint("Some required data was left out:\n\n", err.Message), http.StatusBadRequest)
		return

	case "23503":
		// foreign key violation
		logger.Errorf("foreign key violation: %s", err.Message)
		switch r.Method {
		case "DELETE":
			http.Error(w, fmt.Sprint("This record can’t be deleted because another record refers to it:\n\n", err.Detail), http.StatusConflict)
			return
		}

	case "23505":
		// unique constraint violation
		logger.Errorf("unique constraint violation: %s", err.Message)
		http.Error(w, fmt.Sprint("This record contains duplicated data that conflicts with what is already in the database:\n\n", err.Detail), http.StatusConflict)
		return

	case "23514":
		// check constraint violation
		logger.Errorf("check constraint violation: %s", err.Message)
		http.Error(w, fmt.Sprint("This record contains inconsistent or out-of-range data:\n\n", err.Message), http.StatusConflict)
		return

	default:
		msg := err.Message
		if d := err.Detail; d != "" {
			msg += "\n\n" + d
		}
		if h := err.Hint; h != "" {
			msg += "\n\n" + h
		}
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
