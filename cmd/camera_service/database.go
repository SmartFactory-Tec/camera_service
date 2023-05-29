package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SmartFactory-Tec/camera_service/pkg/migrations"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
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

func connectToDb(config DbConfig, logger *zap.SugaredLogger) *pgxpool.Pool {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Hostname, config.Port, config.Database)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Fatal("error parsing db connection string: %w", err)
	}

	testConnection(pool, logger)

	logger.Infow("connected to database", "name", config.Database)

	return pool
}

func testConnection(conn *pgxpool.Pool, logger *zap.SugaredLogger) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {

		logger.Fatalf("could not connect to database: %s", err)
	}
}

func updateDatabaseSchema(config DbConfig, logger *zap.SugaredLogger) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Hostname, config.Port, config.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatalf("could not connect to database for migrations: %w", err)
	}
	goose.SetBaseFS(migrations.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatalf("could not initialize goose loader: %s", err)
	}

	if err := goose.Up(db, "."); err != nil {
		logger.Fatalf("failed to update database schema: %s", err)
	}

	logger.Infow("updated database schema")
}

func HandlePqError(w http.ResponseWriter, r *http.Request, err *pgconn.PgError, logger *zap.SugaredLogger) {
	// TODO create better errors that do not expose internal database names
	logger = logger.Named("HandlePqError")
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
		case "POST":
			http.Error(w, fmt.Sprint("This resource can't be created because a the referenced id does not exist: \n\n", err.Detail), http.StatusConflict)
			return
		case "PATCH":
			http.Error(w, fmt.Sprint("The resource can't be modified because a referenced id does not exist: \n\n", err.Detail), http.StatusConflict)
			return
		case "DELETE":
			http.Error(w, fmt.Sprint("This resource can’t be deleted because another record refers to it:\n\n", err.Detail), http.StatusConflict)
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
