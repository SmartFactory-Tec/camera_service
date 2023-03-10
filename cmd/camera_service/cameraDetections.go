package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SmartFactory-Tec/camera_service/pkg/dbschema"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
	"path"
	"strconv"
	"time"
)

func makeCreateCameraDetectionHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("CreateCameraDetection")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestBody := struct {
			CameraId          int64      `json:"camera_id"`
			InDirection       int32      `json:"in_direction"`
			OutDirection      int32      `json:"out_direction"`
			Counter           int32      `json:"counter"`
			SocialDistancingV int32      `json:"social_distancing_v"`
			DetectionDate     *time.Time `json:"detection_date"`
		}{}

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&requestBody); err != nil {
			err := fmt.Errorf("error decoding request body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		var (
			id  int64
			err error
		)

		if requestBody.DetectionDate != nil {
			params := dbschema.CreateCameraDetectionWithDateParams{
				CameraID:          requestBody.CameraId,
				InDirection:       requestBody.InDirection,
				OutDirection:      requestBody.OutDirection,
				Counter:           requestBody.Counter,
				SocialDistancingV: requestBody.SocialDistancingV,
				DetectionDate:     *requestBody.DetectionDate,
			}

			id, err = queries.CreateCameraDetectionWithDate(ctx, params)
		} else {
			params := dbschema.CreateCameraDetectionParams{
				CameraID:          requestBody.CameraId,
				InDirection:       requestBody.InDirection,
				OutDirection:      requestBody.OutDirection,
				Counter:           requestBody.Counter,
				SocialDistancingV: requestBody.SocialDistancingV,
			}

			id, err = queries.CreateCameraDetection(ctx, params)
		}

		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				HandlePqError(w, r, err, logger)
				return
			}
			err = fmt.Errorf("error creating camera detection: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", path.Join(r.URL.String(), fmt.Sprintf("/%d", id)))
		w.WriteHeader(http.StatusCreated)
	}
}

func makeGetCameraDetectionsHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetCameraDetections")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cameraDetections, err := queries.GetCameraDetections(ctx)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				HandlePqError(w, r, err, logger)
				return
			}
			err := fmt.Errorf("error getting camera detections: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(cameraDetections)
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

func cameraDetectionCtx(queries *dbschema.Queries, logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	logger = logger.Named("cameraDetectionCtx")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			cameraDetectionId, err := strconv.ParseInt(chi.URLParam(r, "cameraDetectionId"), 10, 64)
			if err != nil {
				logger.Errorf("error parsing camera id: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			cameraDetection, err := queries.GetCameraDetection(ctx, cameraDetectionId)
			if err != nil {
				if err, ok := err.(*pq.Error); ok {
					HandlePqError(w, r, err, logger)
					return
				}
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "camera detection not found", http.StatusNotFound)
					return
				}
				err := fmt.Errorf("error getting camera detection: %w", err)
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "cameraDetection", cameraDetection)))
		})

	}
}

func makeGetCameraDetectionHandler(logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetCameraDetectionHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cameraDetection := ctx.Value("cameraDetection")

		body, err := json.Marshal(cameraDetection)
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

func makeUpdateCameraDetectionHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("UpdateCameraDetection")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cameraDetection := ctx.Value("cameraDetection").(dbschema.CameraDetection)

		dec := json.NewDecoder(r.Body)

		var reqBody struct {
			InDirection       *int32     `json:"in_direction"`
			OutDirection      *int32     `json:"out_direction"`
			Counter           *int32     `json:"counter"`
			SocialDistancingV *int32     `json:"social_distancing_v"`
			DetectionDate     *time.Time `json:"detection_date"`
		}

		if err := dec.Decode(&reqBody); err != nil {
			err = fmt.Errorf("invalid reqBody: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		params := dbschema.UpdateCameraDetectionParams{
			ID:                cameraDetection.ID,
			InDirection:       cameraDetection.InDirection,
			OutDirection:      cameraDetection.OutDirection,
			Counter:           cameraDetection.Counter,
			SocialDistancingV: cameraDetection.SocialDistancingV,
			DetectionDate:     cameraDetection.DetectionDate,
		}

		if reqBody.InDirection != nil {
			params.InDirection = *reqBody.InDirection
		}
		if reqBody.OutDirection != nil {
			params.OutDirection = *reqBody.OutDirection
		}
		if reqBody.Counter != nil {
			params.Counter = *reqBody.Counter
		}
		if reqBody.SocialDistancingV != nil {
			params.SocialDistancingV = *reqBody.SocialDistancingV
		}
		if reqBody.DetectionDate != nil {
			params.DetectionDate = *reqBody.DetectionDate
		}

		cameraDetection, err := queries.UpdateCameraDetection(ctx, params)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				HandlePqError(w, r, err, logger)
				return
			}
			err = fmt.Errorf("error updating camera detection: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resBody, err := json.Marshal(cameraDetection)
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

func makeDeleteCameraDetectionHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("DeleteCameraDetection")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cameraDetection := ctx.Value("cameraDetection").(dbschema.CameraDetection)

		if err := queries.DeleteCameraDetection(ctx, cameraDetection.ID); err != nil {
			if err, ok := err.(*pq.Error); ok {
				HandlePqError(w, r, err, logger)
				return
			}
			err := fmt.Errorf("error deleting camera detection: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}

func makeGetCameraDetectionsByCameraHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetCameraDetectionsByCamera")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		camera := ctx.Value("camera").(dbschema.Camera)

		cameraDetections, err := queries.GetCameraDetectionsFromCamera(ctx, camera.ID)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				HandlePqError(w, r, err, logger)
				return
			}
			err := fmt.Errorf("error getting camera detections: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(cameraDetections)
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
