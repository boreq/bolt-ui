package domain

import (
	"math"

	"github.com/boreq/errors"
)

type Distance struct {
	distance float64
}

func NewDistance(distance float64) (Distance, error) {
	if math.IsNaN(distance) {
		return Distance{}, errors.New("not a number")
	}

	return Distance{
		distance: distance,
	}, nil
}

func MustNewDistance(distance float64) Distance {
	v, err := NewDistance(distance)
	if err != nil {
		panic(err)
	}
	return v
}

func (s Distance) Add(o Distance) Distance {
	return MustNewDistance(s.distance + o.distance)
}

func (s Distance) Float64() float64 {
	return s.distance
}
