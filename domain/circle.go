package domain

import "github.com/boreq/errors"

type Circle struct {
	center Position
	radius float64
}

func NewCircle(center Position, radius float64) (Circle, error) {
	if radius <= 0 {
		return Circle{}, errors.New("radius must be positive")
	}

	return Circle{
		center: center,
		radius: radius,
	}, nil
}

func MustNewCircle(center Position, radius float64) Circle {
	c, err := NewCircle(center, radius)
	if err != nil {
		panic(err)
	}
	return c
}

func (c Circle) Center() Position {
	return c.center
}

func (c Circle) Radius() float64 {
	return c.radius
}

func (c Circle) Contains(p Position) bool {
	return c.center.Distance(p).Float64() <= c.radius
}

func (c Circle) IsZero() bool {
	return c == Circle{}
}
