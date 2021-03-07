package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
)

type position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func toPosition(p domain.Position) position {
	return position{
		Latitude:  p.Latitude().Float64(),
		Longitude: p.Longitude().Float64(),
	}
}

func fromPosition(p position) (domain.Position, error) {
	latitude, err := domain.NewLatitude(p.Latitude)
	if err != nil {
		return domain.Position{}, errors.Wrap(err, "could not create a latitude")
	}

	longitude, err := domain.NewLongitude(p.Longitude)
	if err != nil {
		return domain.Position{}, errors.Wrap(err, "could not create a longitude")
	}

	return domain.NewPosition(latitude, longitude), nil
}

type circle struct {
	Center position `json:"position"`
	Radius float64  `json:"radius"`
}

func toCircle(c domain.Circle) circle {
	return circle{
		Center: toPosition(c.Center()),
		Radius: c.Radius(),
	}
}

func fromCircle(c circle) (domain.Circle, error) {
	p, err := fromPosition(c.Center)
	if err != nil {
		return domain.Circle{}, errors.Wrap(err, "cannot map position")
	}

	return domain.NewCircle(p, c.Radius)
}
