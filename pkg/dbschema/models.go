// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package dbschema

import (
	"database/sql"
	"time"
)

type Camera struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	ConnectionString string `json:"connection_string"`
	LocationText     string `json:"location_text"`
	LocationID       int32  `json:"location_id"`
}

type CameraDetection struct {
	ID                int64         `json:"id"`
	CameraID          sql.NullInt64 `json:"camera_id"`
	InDirection       int32         `json:"in_direction"`
	OutDirection      int32         `json:"out_direction"`
	Counter           int32         `json:"counter"`
	SocialDistancingV int32         `json:"social_distancing_v"`
	DetectionDate     time.Time     `json:"detection_date"`
}

type Location struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
