// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: cameras.sql

package dbschema

import (
	"context"

	"github.com/SmartFactory-Tec/camera_service/pkg/dbenums"
	"github.com/jackc/pgx/v5/pgtype"
)

const createCamera = `-- name: CreateCamera :one
insert into cameras(name, connection_string, location_text, location_id, orientation)
values ($1, $2, $3, $4, $5)
returning id, name, connection_string, location_text, location_id, orientation
`

type CreateCameraParams struct {
	Name             string              `json:"name"`
	ConnectionString string              `json:"connection_string"`
	LocationText     string              `json:"location_text"`
	LocationID       int32               `json:"location_id"`
	Orientation      dbenums.Orientation `json:"orientation"`
}

func (q *Queries) CreateCamera(ctx context.Context, arg CreateCameraParams) (Camera, error) {
	row := q.db.QueryRow(ctx, createCamera,
		arg.Name,
		arg.ConnectionString,
		arg.LocationText,
		arg.LocationID,
		arg.Orientation,
	)
	var i Camera
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ConnectionString,
		&i.LocationText,
		&i.LocationID,
		&i.Orientation,
	)
	return i, err
}

const deleteCamera = `-- name: DeleteCamera :exec
delete
from cameras
where id = $1
`

func (q *Queries) DeleteCamera(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteCamera, id)
	return err
}

const getCamera = `-- name: GetCamera :one
select id, name, connection_string, location_text, location_id, orientation
from cameras
where id = $1
`

func (q *Queries) GetCamera(ctx context.Context, id int64) (Camera, error) {
	row := q.db.QueryRow(ctx, getCamera, id)
	var i Camera
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ConnectionString,
		&i.LocationText,
		&i.LocationID,
		&i.Orientation,
	)
	return i, err
}

const getCameras = `-- name: GetCameras :many
select id, name, connection_string, location_text, location_id, orientation
from cameras
order by id
`

func (q *Queries) GetCameras(ctx context.Context) ([]Camera, error) {
	rows, err := q.db.Query(ctx, getCameras)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Camera{}
	for rows.Next() {
		var i Camera
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ConnectionString,
			&i.LocationText,
			&i.LocationID,
			&i.Orientation,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCamera = `-- name: UpdateCamera :one
update cameras
set name              = coalesce($2, name),
    connection_string = coalesce($3, connection_string),
    location_text     = coalesce($4, location_text),
    location_id       = coalesce($5, location_id),
    orientation       = coalesce($6, orientation)
where id = $1
returning id, name, connection_string, location_text, location_id, orientation
`

type UpdateCameraParams struct {
	ID               int64                   `json:"id"`
	Name             pgtype.Text             `json:"name"`
	ConnectionString pgtype.Text             `json:"connection_string"`
	LocationText     pgtype.Text             `json:"location_text"`
	LocationID       pgtype.Int4             `json:"location_id"`
	Orientation      dbenums.NullOrientation `json:"orientation"`
}

func (q *Queries) UpdateCamera(ctx context.Context, arg UpdateCameraParams) (Camera, error) {
	row := q.db.QueryRow(ctx, updateCamera,
		arg.ID,
		arg.Name,
		arg.ConnectionString,
		arg.LocationText,
		arg.LocationID,
		arg.Orientation,
	)
	var i Camera
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ConnectionString,
		&i.LocationText,
		&i.LocationID,
		&i.Orientation,
	)
	return i, err
}
