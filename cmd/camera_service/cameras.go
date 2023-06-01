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

func cameraCtx(queries *dbschema.Queries, logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	logger = logger.Named("cameraCtx")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			cameraId, err := strconv.ParseInt(chi.URLParam(r, "cameraId"), 10, 64)
			if err != nil {
				err := fmt.Errorf("error parsing camera id: %w", err)
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			camera, err := queries.GetCamera(ctx, cameraId)

			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				HandlePqError(w, r, pgErr, logger)
			} else if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "camera not found", http.StatusNotFound)
			} else if err != nil {
				err := fmt.Errorf("error getting camera: %w", err)
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "camera", camera)))
			}
		})
	}
}

func getCameras(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetCameras")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cameras, err := queries.GetCameras(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
		} else if err != nil {
			err := fmt.Errorf("error getting cameras: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(cameras)
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

func getCamera(logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetCamera")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		camera := ctx.Value("camera")

		body, err := json.Marshal(camera)
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

func postCamera(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("CreateCamera")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var params dbschema.CreateCameraParams

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&params); err != nil {
			err := fmt.Errorf("error decoding request body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		camera, err := queries.CreateCamera(ctx, params)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			err = fmt.Errorf("error creating camera detection: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(camera)
		if err != nil {
			err = fmt.Errorf("error marshaling body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", path.Join(r.URL.String(), fmt.Sprintf("/%d", camera.ID)))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(body); err != nil {
			err = fmt.Errorf("error writing body: %w", err)
			logger.Error(err)
		}
	}
}

func patchCamera(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("patchCamera")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		camera := ctx.Value("camera").(dbschema.Camera)

		dec := json.NewDecoder(r.Body)

		var params dbschema.UpdateCameraParams

		if err := dec.Decode(&params); err != nil {
			err = fmt.Errorf("invalid body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		params.ID = camera.ID

		camera, err := queries.UpdateCamera(ctx, params)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			err = fmt.Errorf("error updating camera: %s", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resBody, err := json.Marshal(camera)
		if err != nil {
			err := fmt.Errorf("error marshaling json body: %s", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(resBody); err != nil {
			logger.Errorf("error writing json body: %s", err)
		}
	}
}

func deleteCamera(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("DeleteCamera")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		camera := ctx.Value("camera").(dbschema.Camera)

		err := queries.DeleteCamera(ctx, camera.ID)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
		} else if err != nil {
			err := fmt.Errorf("error deleting camera: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
