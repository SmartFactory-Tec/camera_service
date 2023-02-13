// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: camera_detections.sql

package dbschema

import (
	"context"
	"time"
)

const createCameraDetection = `-- name: CreateCameraDetection :exec
insert into camera_detections(camera_id, in_direction, out_direction, counter, social_distancing_v, detection_date)
values ($1, $2, $3, $4, $5, $6)
`

type CreateCameraDetectionParams struct {
	CameraID          int64
	InDirection       int32
	OutDirection      int32
	Counter           int32
	SocialDistancingV int32
	DetectionDate     time.Time
}

func (q *Queries) CreateCameraDetection(ctx context.Context, arg CreateCameraDetectionParams) error {
	_, err := q.db.ExecContext(ctx, createCameraDetection,
		arg.CameraID,
		arg.InDirection,
		arg.OutDirection,
		arg.Counter,
		arg.SocialDistancingV,
		arg.DetectionDate,
	)
	return err
}

const deleteCameraDetection = `-- name: DeleteCameraDetection :exec
delete
from camera_detections
where id = $1
`

func (q *Queries) DeleteCameraDetection(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteCameraDetection, id)
	return err
}

const getCameraDetection = `-- name: GetCameraDetection :one
select id, camera_id, in_direction, out_direction, counter, social_distancing_v, detection_date
from camera_detections
where id = $1
`

func (q *Queries) GetCameraDetection(ctx context.Context, id int64) (CameraDetection, error) {
	row := q.db.QueryRowContext(ctx, getCameraDetection, id)
	var i CameraDetection
	err := row.Scan(
		&i.ID,
		&i.CameraID,
		&i.InDirection,
		&i.OutDirection,
		&i.Counter,
		&i.SocialDistancingV,
		&i.DetectionDate,
	)
	return i, err
}

const getCameraDetections = `-- name: GetCameraDetections :many
select id, camera_id, in_direction, out_direction, counter, social_distancing_v, detection_date
from camera_detections
order by id
`

func (q *Queries) GetCameraDetections(ctx context.Context) ([]CameraDetection, error) {
	rows, err := q.db.QueryContext(ctx, getCameraDetections)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CameraDetection
	for rows.Next() {
		var i CameraDetection
		if err := rows.Scan(
			&i.ID,
			&i.CameraID,
			&i.InDirection,
			&i.OutDirection,
			&i.Counter,
			&i.SocialDistancingV,
			&i.DetectionDate,
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

const updateCameraDetection = `-- name: UpdateCameraDetection :one
update camera_detections
set in_direction        = $2,
    out_direction       = $3,
    counter             = $4,
    social_distancing_v = $5,
    detection_date      = $6
where id = $1
returning id, camera_id, in_direction, out_direction, counter, social_distancing_v, detection_date
`

type UpdateCameraDetectionParams struct {
	ID                int64
	InDirection       int32
	OutDirection      int32
	Counter           int32
	SocialDistancingV int32
	DetectionDate     time.Time
}

func (q *Queries) UpdateCameraDetection(ctx context.Context, arg UpdateCameraDetectionParams) (CameraDetection, error) {
	row := q.db.QueryRowContext(ctx, updateCameraDetection,
		arg.ID,
		arg.InDirection,
		arg.OutDirection,
		arg.Counter,
		arg.SocialDistancingV,
		arg.DetectionDate,
	)
	var i CameraDetection
	err := row.Scan(
		&i.ID,
		&i.CameraID,
		&i.InDirection,
		&i.OutDirection,
		&i.Counter,
		&i.SocialDistancingV,
		&i.DetectionDate,
	)
	return i, err
}
