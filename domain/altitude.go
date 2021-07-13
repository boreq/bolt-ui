package domain

import (
	"math"

	"github.com/boreq/errors"
)

type Altitude struct {
	altitude float64
}

func NewAltitude(altitude float64) (Altitude, error) {
	if math.IsNaN(altitude) {
		return Altitude{}, errors.New("not a number")
	}

	return Altitude{
		altitude: altitude,
	}, nil
}

func MustNewAltitude(altitude float64) Altitude {
	a, err := NewAltitude(altitude)
	if err != nil {
		panic(err)
	}
	return a
}

func (a Altitude) Float64() float64 {
	return a.altitude
}
