package domain

import (
	"math"
	"time"

	"github.com/boreq/errors"
)

type Speed struct {
	speed float64
}

func NewSpeed(speed float64) (Speed, error) {
	if math.IsNaN(speed) {
		return Speed{}, errors.New("not a number")
	}

	return Speed{
		speed: speed,
	}, nil
}

func NewSpeedFromDistanceAndDuration(distance Distance, duration time.Duration) (Speed, error) {
	return NewSpeed(distance.Float64() / duration.Seconds())
}

func MustNewSpeed(speed float64) Speed {
	v, err := NewSpeed(speed)
	if err != nil {
		panic(err)
	}
	return v
}

func (s Speed) Float64() float64 {
	return s.speed
}

func (s Speed) IsZero() bool {
	return s == Speed{}
}
