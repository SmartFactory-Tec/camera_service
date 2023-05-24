package dbenums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Orientation string

const (
	CameraOrientationVertical           Orientation = "vertical"
	CameraOrientationHorizontal         Orientation = "horizontal"
	CameraOrientationInvertedVertical   Orientation = "inverted_vertical"
	CameraOrientationInvertedHorizontal Orientation = "inverted_horizontal"
)

func (co Orientation) Value() (driver.Value, error) {
	return string(co), nil
}

func (co *Orientation) Scan(src any) error {
	switch src := src.(type) {
	case string:
		if src == string(CameraOrientationHorizontal) ||
			src == string(CameraOrientationVertical) ||
			src == string(CameraOrientationInvertedHorizontal) ||
			src == string(CameraOrientationInvertedVertical) {
			*co = Orientation(src)
			return nil
		} else {
			return fmt.Errorf("invalid value for enum CameraOrientation")
		}
	default:
		return fmt.Errorf("invalid value for enum CameraOrientation")
	}
}

func (co *Orientation) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*co))
}

func (co *Orientation) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	if s == string(CameraOrientationHorizontal) ||
		s == string(CameraOrientationVertical) ||
		s == string(CameraOrientationInvertedHorizontal) ||
		s == string(CameraOrientationInvertedVertical) {
		*co = Orientation(s)
		return nil
	} else {
		return fmt.Errorf("invalid value for enum CameraOrientation")
	}
}

type NullOrientation struct {
	Orientation Orientation
	Valid       bool
}

func (co NullOrientation) Value() (driver.Value, error) {
	if !co.Valid {
		return nil, nil
	}
	return co.Orientation.Value()
}

func (co *NullOrientation) Scan(src any) error {
	if err := co.Orientation.Scan(src); err != nil {
		return err
	}
	co.Valid = true
	return nil
}

func (co *NullOrientation) MarshalJSON() ([]byte, error) {
	if !co.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(co.Orientation)
}

func (co *NullOrientation) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &co.Orientation); err != nil {
		return err
	}
	co.Valid = true
	return nil
}
