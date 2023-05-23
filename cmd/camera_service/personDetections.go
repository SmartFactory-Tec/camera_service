package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SmartFactory-Tec/camera_service/pkg/dbschema"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"net/http"
	"path"
	"strconv"
)

func makeCreatePersonDetectionHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("CreatePersonDetection")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestBody := dbschema.CreatePersonDetectionParams{}

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&requestBody); err != nil {
			err := fmt.Errorf("error decoding request body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		personDetection, err := queries.CreatePersonDetection(ctx, requestBody)

		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			HandlePqError(w, r, pqErr, logger)
			return
		} else if err != nil {
			err = fmt.Errorf("error creating camera detection: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(personDetection)
		if err != nil {
			err = fmt.Errorf("error marshaling camera detection: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", path.Join(r.URL.String(), fmt.Sprintf("/%d", personDetection.ID)))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err = w.Write(body); err != nil {
			err = fmt.Errorf("error writing body: %w", err)
			logger.Error(err)
		}
	}
}

func makeGetPersonDetectionsHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetPersonDetections")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		offset, err := strconv.ParseInt(chi.URLParam(r, "offset"), 10, 32)
		if err != nil {
			err := fmt.Errorf("request does not contain parameter offset: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		count, err := strconv.ParseInt(chi.URLParam(r, "count"), 10, 32)
		if err != nil {
			err := fmt.Errorf("request does not contain parameter count: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		params := dbschema.GetPersonDetectionsParams{
			DetectionOffset: int32(offset),
			Count:           int32(count),
		}
		personDetections, err := queries.GetPersonDetections(ctx, params)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			err := fmt.Errorf("error getting person detections: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(personDetections)
		if err != nil {
			logger.Errorf("error marshaling json body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(body); err != nil {
			logger.Errorf("error writing json body: %s", err)
		}
	}
}

func personDetectionCtx(queries *dbschema.Queries, logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	logger = logger.Named("personDetectionCtx")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			personDetectionId, err := strconv.ParseInt(chi.URLParam(r, "personDetectionId"), 10, 64)
			if err != nil {
				logger.Errorf("error parsing person detection id: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			personDetection, err := queries.GetPersonDetection(ctx, personDetectionId)
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				HandlePqError(w, r, pgErr, logger)
				return
			} else if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "person detection not found", http.StatusNotFound)
				return
			} else if err != nil {
				err := fmt.Errorf("error getting person detection: %w", err)
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "personDetection", personDetection)))
		})

	}
}

func makeGetPersonDetectionHandler(logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetPersonDetectionHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		personDetection := ctx.Value("personDetection")

		body, err := json.Marshal(personDetection)
		if err != nil {
			err := fmt.Errorf("error marshaling json body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if _, err := w.Write(body); err != nil {
			logger.Errorf("error writing json body: %s", err)
		}
	}
}

func makeUpdatePersonDetectionHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("UpdatePersonDetection")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		personDetection := ctx.Value("personDetection").(dbschema.PersonDetection)

		dec := json.NewDecoder(r.Body)

		params := dbschema.UpdatePersonDetectionParams{}

		if err := dec.Decode(&params); err != nil {
			err = fmt.Errorf("invalid request body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		personDetection, err := queries.UpdatePersonDetection(ctx, params)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			err = fmt.Errorf("error updating person detection: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(personDetection)
		if err != nil {
			err := fmt.Errorf("error marshaling json body: %s", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(body); err != nil {
			logger.Errorf("error writing json body: %s", err)
		}
	}
}

func makeDeletePersonDetectionHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("DeletePersonDetection")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		personDetection := ctx.Value("personDetection").(dbschema.PersonDetection)

		err := queries.DeletePersonDetection(ctx, personDetection.ID)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			err := fmt.Errorf("error deleting person detection: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func makeGetPersonDetectionsByCameraHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetPersonDetectionsByCamera")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		camera := ctx.Value("camera").(dbschema.Camera)
		offset, err := strconv.ParseInt(chi.URLParam(r, "offset"), 10, 32)
		if err != nil {
			err := fmt.Errorf("request does not contain parameter offset: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		count, err := strconv.ParseInt(chi.URLParam(r, "count"), 10, 32)
		if err != nil {
			err := fmt.Errorf("request does not contain parameter count: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		params := dbschema.GetPersonDetectionsForCameraParams{
			CameraID:        camera.ID,
			DetectionOffset: int32(offset),
			Count:           int32(count),
		}
		personDetections, err := queries.GetPersonDetectionsForCamera(ctx, params)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			err := fmt.Errorf("error getting person detections: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(personDetections)
		if err != nil {
			logger.Errorf("error marshaling json body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(body); err != nil {
			logger.Errorf("error writing json body: %s", err)
		}
	}
}
