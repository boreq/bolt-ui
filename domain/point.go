package domain

import (
	"time"

	"github.com/boreq/errors"
)

type Point struct {
	time     time.Time
	position Position
	altitude Altitude
}

func NewPoint(t time.Time, position Position, altitude Altitude) (Point, error) {
	if t.IsZero() {
		return Point{}, errors.New("zero value of time")
	}

	return Point{
		time:     t,
		position: position,
		altitude: altitude,
	}, nil
}

func MustNewPoint(t time.Time, position Position, altitude Altitude) Point {
	v, err := NewPoint(t, position, altitude)
	if err != nil {
		panic(err)
	}
	return v
}

func (p Point) Time() time.Time {
	return p.time
}
