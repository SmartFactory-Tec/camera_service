// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package dbschema

import (
	"time"
)

type Camera struct {
	ID               int64
	Name             string
	ConnectionString string
	LocationText     string
	LocationID       int32
}

type CameraDetection struct {
	ID                int64
	CameraID          int64
	InDirection       int32
	OutDirection      int32
	Counter           int32
	SocialDistancingV int32
	DetectionDate     time.Time
}

type Location struct {
	ID          int64
	Name        string
	Description string
}
