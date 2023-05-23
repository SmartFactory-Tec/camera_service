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

func makeCreateLocationHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("CreateLocation")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		params := dbschema.CreateLocationParams{}

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&params); err != nil {
			err := fmt.Errorf("error decoding request body: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := queries.CreateLocation(ctx, params)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			err = fmt.Errorf("error creating location: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", path.Join(r.URL.String(), fmt.Sprintf("/%d", id)))
		w.WriteHeader(http.StatusCreated)
	}
}

func makeGetLocationsHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetLocations")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		locations, err := queries.GetLocations(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			logger.Error(fmt.Errorf("error getting locations from db: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := json.Marshal(locations)
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

func locationCtx(queries *dbschema.Queries, logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	logger = logger.Named("locationCtx")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			locationId, err := strconv.ParseInt(chi.URLParam(r, "locationId"), 10, 64)
			if err != nil {
				logger.Errorf("error parsing location id: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			location, err := queries.GetLocation(ctx, locationId)
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				HandlePqError(w, r, pgErr, logger)
				return
			} else if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "location not found", http.StatusNotFound)
				return
			} else if err != nil {
				err := fmt.Errorf("error getting location: %w", err)
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "location", location)))
		})

	}
}

func makeGetLocationHandler(logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("GetLocation")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		location := ctx.Value("location")

		body, err := json.Marshal(location)
		if err != nil {
			logger.Errorf("error marshaling json body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if _, err := w.Write(body); err != nil {
			logger.Errorf("error writing json body: %s", err)
		}
	}
}

func makeUpdateLocationHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("UpdateLocation")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		location := ctx.Value("location").(dbschema.Location)

		dec := json.NewDecoder(r.Body)

		var reqBody struct {
			Name        *string `json:"name"`
			Description *string `json:"description"`
		}

		if err := dec.Decode(&reqBody); err != nil {
			err = fmt.Errorf("invalid reqBody: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		params := dbschema.UpdateLocationParams{
			ID:          location.ID,
			Name:        location.Name,
			Description: location.Description,
		}

		if reqBody.Name != nil {
			params.Name = *reqBody.Name
		}
		if reqBody.Description != nil {
			params.Description = *reqBody.Description
		}

		location, err := queries.UpdateLocation(ctx, params)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			err = fmt.Errorf("error updating location: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resBody, err := json.Marshal(location)
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

func makeDeleteLocationHandler(queries *dbschema.Queries, logger *zap.SugaredLogger) http.HandlerFunc {
	logger = logger.Named("DeleteLocation")
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		location := ctx.Value("location").(dbschema.Location)

		err := queries.DeleteLocation(ctx, location.ID)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandlePqError(w, r, pgErr, logger)
			return
		} else if err != nil {
			err := fmt.Errorf("error deleting location: %w", err)
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
