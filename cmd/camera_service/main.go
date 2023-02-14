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

	r.Route("/cameras", func(r chi.Router) {
		r.Get("/", makeGetCamerasHandler(queries, logger))
		r.Post("/", makeCreateCameraHandler(queries, logger))

		r.Route("/{cameraId}", func(r chi.Router) {
			r.Use(cameraCtx(queries, logger))
			r.Get("/", makeGetCameraHandler(logger))
			r.Patch("/", makeUpdateCameraHandler(queries, logger))
			r.Delete("/", makeDeleteCameraHandler(queries, logger))

			r.Get("/{cameraId}/cameraDetections", makeGetCameraDetectionsByCameraHandler(queries, logger))
		})

	})

	r.Route("/cameraDetections", func(r chi.Router) {
		r.Get("/", makeGetCameraDetectionsHandler(queries, logger))
		r.Get("/", makeCreateCameraDetectionHandler(queries, logger))

		r.Route("/{cameraDetectionId}", func(r chi.Router) {
			r.Use(cameraDetectionCtx(queries, logger))
			r.Get("/", makeGetCameraDetectionHandler(logger))
			r.Patch("/", makeUpdateCameraDetectionHandler(queries, logger))
			r.Delete("/", makeDeleteCameraDetectionHandler(queries, logger))
		})
	})

	logger.Infof("starting server on port %d", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
	if err != nil {
		logger.Fatal(fmt.Errorf("http server error: %w", err))
	}

}
