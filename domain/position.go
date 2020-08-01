package domain

import "github.com/boreq/errors"

type Position struct {
	latitude  Latitude
	longitude Longitude
}

func NewPosition(latitude Latitude, longitude Longitude) Position {
	return Position{
		latitude:  latitude,
		longitude: longitude,
	}
}

func (p Position) Latitude() Latitude {
	return p.latitude
}

func (p Position) Longitude() Longitude {
	return p.longitude
}

type Longitude struct {
	longitude float64
}

func NewLongitude(longitude float64) (Longitude, error) {
	if longitude < -180 || longitude > 180 {
		return Longitude{}, errors.New("invalid longitude")
	}

	return Longitude{
		longitude: longitude,
	}, nil
}

func MustNewLongitude(longitude float64) Longitude {
	v, err := NewLongitude(longitude)
	if err != nil {
		panic(err)
	}
	return v
}

func (l Longitude) Float64() float64 {
	return l.longitude
}

type Latitude struct {
	latitude float64
}

func NewLatitude(latitude float64) (Latitude, error) {
	if latitude < -90 || latitude > 90 {
		return Latitude{}, errors.New("invalid latitude")
	}

	return Latitude{
		latitude: latitude,
	}, nil
}

func MustNewLatitude(latitude float64) Latitude {
	v, err := NewLatitude(latitude)
	if err != nil {
		panic(err)
	}
	return v
}

func (l Latitude) Float64() float64 {
	return l.latitude
}
