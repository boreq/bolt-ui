package domain

import (
	"math"

	"github.com/boreq/errors"
)

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

func (p Position) Distance(o Position) Distance {
	var la1, lo1, la2, lo2, r float64
	la1 = p.latitude.Float64() * math.Pi / 180
	lo1 = p.longitude.Float64() * math.Pi / 180
	la2 = o.latitude.Float64() * math.Pi / 180
	lo2 = o.longitude.Float64() * math.Pi / 180
	r = 6378100 // Earth radius in meters
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return MustNewDistance(2 * r * math.Asin(math.Sqrt(h)))
}

type Longitude struct {
	longitude float64
}

func NewLongitude(longitude float64) (Longitude, error) {
	if longitude < -180 || longitude > 180 {
		return Longitude{}, errors.New("invalid longitude")
	}

	if math.IsNaN(longitude) {
		return Longitude{}, errors.New("not a number")
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

	if math.IsNaN(latitude) {
		return Latitude{}, errors.New("not a number")
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

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
