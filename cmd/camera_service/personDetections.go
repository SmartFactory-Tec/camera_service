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
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"net/http"
	"path"
	"strconv"
)

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

func getPersonDetections(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("getPersonDetections")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		offsetStr := r.URL.Query().Get("offset")
		countStr := r.URL.Query().Get("count")
		offset, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			offset = 0
		}
		count, err := strconv.ParseInt(countStr, 10, 32)
		if err != nil {
			err := fmt.Errorf("request does not contain required parameter count: %w", err)
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

func getPersonDetection(logger *zap.SugaredLogger) http.HandlerFunc {
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

func postPersonDetection(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("postPersonDetection")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		dec := json.NewDecoder(r.Body)

		var params dbschema.CreatePersonDetectionParams
		if err := dec.Decode(&params); err != nil {
			err := fmt.Errorf("error decoding request body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		personDetection, err := queries.CreatePersonDetection(ctx, params)

		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			HandlePqError(w, r, pqErr, logger)
			return
		} else if err != nil {
			err = fmt.Errorf("error creating person detection: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(personDetection)
		if err != nil {
			err = fmt.Errorf("error marshaling person detection: %w", err)
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

func getDailyPersonDetectionsCount(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("getDailyPersonDetectionsCount")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		camera := ctx.Value("camera").(dbschema.Camera)

		daysStr := r.URL.Query().Get("days")
		monthsStr := r.URL.Query().Get("months")
		days, err := strconv.ParseInt(daysStr, 10, 32)
		if err != nil {
			days = 0
		}

		months, err := strconv.ParseInt(monthsStr, 10, 32)
		if err != nil {
			months = 0
		}

		params := dbschema.GetDailyPersonDetectionsCountParams{
			CameraID: camera.ID,
			Interval: pgtype.Interval{
				Microseconds: 0,
				Days:         int32(days),
				Months:       int32(months),
				Valid:        true,
			},
		}

		dailyPersonDetectionsCount, err := queries.GetDailyPersonDetectionsCount(ctx, params)

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

		body, err := json.Marshal(dailyPersonDetectionsCount)
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

func patchPersonDetection(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("UpdatePersonDetection")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		personDetection := ctx.Value("personDetection").(dbschema.PersonDetection)

		dec := json.NewDecoder(r.Body)

		var params dbschema.UpdatePersonDetectionParams

		if err := dec.Decode(&params); err != nil {
			err = fmt.Errorf("invalid request body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		params.ID = personDetection.ID

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

func deletePersonDetection(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
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

func getCameraPersonDetections(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetPersonDetectionsByCamera")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		camera := ctx.Value("camera").(dbschema.Camera)

		offsetStr := r.URL.Query().Get("offset")
		countStr := r.URL.Query().Get("count")
		offset, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			offset = 0
		}
		count, err := strconv.ParseInt(countStr, 10, 32)
		if err != nil {
			err := fmt.Errorf("request does not contain required parameter count: %w", err)
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

func postCameraPersonDetection(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("postPersonDetection")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		camera := ctx.Value("camera").(dbschema.Camera)

		dec := json.NewDecoder(r.Body)

		var params dbschema.CreatePersonDetectionParams
		if err := dec.Decode(&params); err != nil {
			err := fmt.Errorf("error decoding request body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		params.CameraID = camera.ID

		personDetection, err := queries.CreatePersonDetection(ctx, params)

		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			HandlePqError(w, r, pqErr, logger)
			return
		} else if err != nil {
			err = fmt.Errorf("error creating person detection: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(personDetection)
		if err != nil {
			err = fmt.Errorf("error marshaling person detection: %w", err)
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
