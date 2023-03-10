// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: cameras.sql

package dbschema

import (
	"context"
)

const createCamera = `-- name: CreateCamera :one
insert into cameras(name, connection_string, location_text, location_id)
values ($1, $2, $3, $4)
returning id
`

type CreateCameraParams struct {
	Name             string `json:"name"`
	ConnectionString string `json:"connection_string"`
	LocationText     string `json:"location_text"`
	LocationID       int32  `json:"location_id"`
}

func (q *Queries) CreateCamera(ctx context.Context, arg CreateCameraParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createCamera,
		arg.Name,
		arg.ConnectionString,
		arg.LocationText,
		arg.LocationID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteCamera = `-- name: DeleteCamera :exec
delete
from cameras
where id = $1
`

func (q *Queries) DeleteCamera(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteCamera, id)
	return err
}

const getCamera = `-- name: GetCamera :one
select id, name, connection_string, location_text, location_id
from cameras
where id = $1
`

func (q *Queries) GetCamera(ctx context.Context, id int64) (Camera, error) {
	row := q.db.QueryRowContext(ctx, getCamera, id)
	var i Camera
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ConnectionString,
		&i.LocationText,
		&i.LocationID,
	)
	return i, err
}

const getCameras = `-- name: GetCameras :many
select id, name, connection_string, location_text, location_id
from cameras
order by id
`

func (q *Queries) GetCameras(ctx context.Context) ([]Camera, error) {
	rows, err := q.db.QueryContext(ctx, getCameras)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Camera
	for rows.Next() {
		var i Camera
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ConnectionString,
			&i.LocationText,
			&i.LocationID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCamera = `-- name: UpdateCamera :one
update cameras
set name              = $2,
    connection_string = $3,
    location_text     = $4,
    location_id       = $5
where id = $1
returning id, name, connection_string, location_text, location_id
`

type UpdateCameraParams struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	ConnectionString string `json:"connection_string"`
	LocationText     string `json:"location_text"`
	LocationID       int32  `json:"location_id"`
}

func (q *Queries) UpdateCamera(ctx context.Context, arg UpdateCameraParams) (Camera, error) {
	row := q.db.QueryRowContext(ctx, updateCamera,
		arg.ID,
		arg.Name,
		arg.ConnectionString,
		arg.LocationText,
		arg.LocationID,
	)
	var i Camera
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ConnectionString,
		&i.LocationText,
		&i.LocationID,
	)
	return i, err
}
