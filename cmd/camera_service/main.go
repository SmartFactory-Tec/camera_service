package main

import (
	"fmt"
	"github.com/SmartFactory-Tec/camera_service/pkg/dbschema"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	logger := setupLogger()

	config := loadConfig(logger)
	dbConfig := config.Db

	db := connectToDb(dbConfig, logger)
	updateDatabaseSchema(db, logger)
	queries := dbschema.New(db)

	r := chi.NewRouter()
	r.Use(LogRequests(logger))

	r.Route("/locations", func(r chi.Router) {
		r.Get("/", makeGetLocationsHandler(queries, logger))
		r.Post("/", makeCreateLocationHandler(queries, logger))

		r.Route("/{locationId}", func(r chi.Router) {
			r.Use(locationCtx(queries, logger))
			r.Get("/", makeGetLocationHandler(logger))
			r.Patch("/", makeUpdateLocationHandler(queries, logger))
			r.Delete("/", makeDeleteLocationHandler(queries, logger))
		})
	})

	logger.Infof("starting server on port %d", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
	if err != nil {
		logger.Fatal(fmt.Errorf("http server error: %w", err))
	}

}
