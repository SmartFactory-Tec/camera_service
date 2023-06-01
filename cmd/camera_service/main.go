package main

import (
	"fmt"
	"github.com/SmartFactory-Tec/camera_service/pkg/dbschema"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func main() {
	logger := setupLogger()

	config := loadConfig(logger)
	dbConfig := config.Db

	db := connectToDb(dbConfig, logger)
	updateDatabaseSchema(dbConfig, logger)
	queries := dbschema.New(db)

	var allowedOrigins []string

	if !config.Cors.AllowAllOrigins {
		allowedOrigins = config.Cors.AllowedOrigins
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "OPTIONS", "POST", "PATCH"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"*"},
	}))

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
		r.Get("/", getCameras(queries, logger))
		r.Post("/", postCamera(queries, logger))

		r.Route("/{cameraId}", func(r chi.Router) {
			r.Use(cameraCtx(queries, logger))
			r.Get("/", getCamera(logger))
			r.Patch("/", patchCamera(queries, logger))
			r.Delete("/", deleteCamera(queries, logger))

			r.Get("/personDetections", getCameraPersonDetections(queries, logger))
			r.Post("/personDetections", postCameraPersonDetection(queries, logger)
		})

	})

	r.Route("/personDetections", func(r chi.Router) {
		r.Get("/", getPersonDetections(queries, logger))
		r.Post("/", postPersonDetection(queries, logger))

		r.Route("/{personDetectionId}", func(r chi.Router) {
			r.Use(personDetectionCtx(queries, logger))
			r.Get("/", getPersonDetection(logger))
			r.Patch("/", patchPersonDetection(queries, logger))
			r.Delete("/", deletePersonDetection(queries, logger))
		})
	})

	logger.Infof("starting server on port %d", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
	if err != nil {
		logger.Fatal(fmt.Errorf("http server error: %w", err))
	}

}
