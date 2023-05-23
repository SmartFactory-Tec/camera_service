package dbenums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type CameraOrientation string

const (
	CameraOrientationVertical           CameraOrientation = "vertical"
	CameraOrientationHorizontal         CameraOrientation = "horizontal"
	CameraOrientationInvertedVertical   CameraOrientation = "inverted_vertical"
	CameraOrientationInvertedHorizontal CameraOrientation = "inverted_horizontal"
)

func (co *CameraOrientation) Value() (driver.Value, error) {
	return string(*co), nil
}

func (co *CameraOrientation) Scan(src any) error {
	switch src := src.(type) {
	case string:
		if src == string(CameraOrientationHorizontal) ||
			src == string(CameraOrientationVertical) ||
			src == string(CameraOrientationInvertedHorizontal) ||
			src == string(CameraOrientationInvertedVertical) {
			*co = CameraOrientation(src)
			return nil
		} else {
			return fmt.Errorf("invalid value for enum CameraOrientation")
		}
	default:
		return fmt.Errorf("invalid value for enum CameraOrientation")
	}
}

func (co *CameraOrientation) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*co))
}

func (co *CameraOrientation) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	if s == string(CameraOrientationHorizontal) ||
		s == string(CameraOrientationVertical) ||
		s == string(CameraOrientationInvertedHorizontal) ||
		s == string(CameraOrientationInvertedVertical) {
		*co = CameraOrientation(s)
		return nil
	} else {
		return fmt.Errorf("invalid value for enum CameraOrientation")
	}
}

type NullCameraOrientation struct {
	CameraOrientation CameraOrientation
	Valid             bool
}

func (co *NullCameraOrientation) Value() (driver.Value, error) {
	if !co.Valid {
		return nil, nil
	}
	return co.CameraOrientation.Value()
}

func (co *NullCameraOrientation) Scan(src any) error {
	if err := co.CameraOrientation.Scan(src); err != nil {
		return err
	}
	co.Valid = true
	return nil
}

func (co *NullCameraOrientation) MarshalJSON() ([]byte, error) {
	if !co.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(co.CameraOrientation)
}

func (co *NullCameraOrientation) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &co.CameraOrientation); err != nil {
		return err
	}
	co.Valid = true
	return nil
}
