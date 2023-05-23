package dbenums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Direction string

const (
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
	DirectionNone  Direction = "none"
)

func (d *Direction) Value() (driver.Value, error) {
	return string(*d), nil
}

func (d *Direction) Scan(src any) error {
	switch src := src.(type) {
	case string:
		if src == string(DirectionNone) ||
			src == string(DirectionLeft) ||
			src == string(DirectionRight) {
			*d = Direction(src)
			return nil
		} else {
			return fmt.Errorf("invalid value for enum Direction")
		}
	default:
		return fmt.Errorf("invalid value for enum Direction")
	}
}

func (d *Direction) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*d))
}

func (d *Direction) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	if s == string(DirectionNone) ||
		s == string(DirectionLeft) ||
		s == string(DirectionRight) {
		*d = Direction(s)
		return nil
	} else {
		return fmt.Errorf("invalid value for enum Direction")
	}
}

type NullDirection struct {
	Direction Direction
	Valid     bool
}

func (d *NullDirection) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Direction.Value()
}

func (d *NullDirection) Scan(src any) error {
	if err := d.Direction.Scan(src); err != nil {
		return err
	}
	d.Valid = true
	return nil
}

func (d *NullDirection) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(d.Direction)
}

func (d *NullDirection) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &d.Direction); err != nil {
		return err
	}
	d.Valid = true
	return nil
}
