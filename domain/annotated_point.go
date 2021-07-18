package domain

import (
	"time"

	"github.com/boreq/errors"
)

type AnnotatedPoint struct {
	point              Point
	speed              Speed
	cumulativeDistance Distance
}

func NewAnnotatedPoint(point Point, speed Speed, cumulativeDistance Distance) (AnnotatedPoint, error) {
	if point.IsZero() {
		return AnnotatedPoint{}, errors.New("zero value of point")
	}

	return AnnotatedPoint{
		point:              point,
		speed:              speed,
		cumulativeDistance: cumulativeDistance,
	}, nil
}

func MustNewAnnotatedPoint(point Point, speed Speed, cumulativeDistance Distance) AnnotatedPoint {
	v, err := NewAnnotatedPoint(point, speed, cumulativeDistance)
	if err != nil {
		panic(err)
	}
	return v
}

func (p AnnotatedPoint) Point() Point {
	return p.point
}

func (p AnnotatedPoint) Time() time.Time {
	return p.point.time
}

func (p AnnotatedPoint) Position() Position {
	return p.point.position
}

func (p AnnotatedPoint) Altitude() Altitude {
	return p.point.altitude
}

func (p *AnnotatedPoint) SetSpeedToZero() {
	p.speed = MustNewSpeed(0)
}

func (p AnnotatedPoint) Speed() Speed {
	return p.speed
}

func (p AnnotatedPoint) CumulativeDistance() Distance {
	return p.cumulativeDistance
}
