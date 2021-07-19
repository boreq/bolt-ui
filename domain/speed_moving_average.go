package domain

import (
	"time"
)

type SpeedMovingAverage struct {
	timePeriod time.Duration
	points     []Point
}

func NewSpeedMovingAverage(timePeriod time.Duration) *SpeedMovingAverage {
	return &SpeedMovingAverage{
		timePeriod: timePeriod,
	}
}

func (a *SpeedMovingAverage) AddPoint(p Point) {
	a.points = append(a.points, p)

	for {
		if len(a.points) <= 2 {
			break
		}

		if a.points[0].Time().Add(a.timePeriod).After(p.time) {
			break
		}

		a.points = append(a.points[:0], a.points[1:]...)
	}
}

func (a *SpeedMovingAverage) Speed() (Speed, error) {
	if len(a.points) < 2 {
		return NewSpeed(0)
	}

	distance := a.distanceBetweenPoints()
	duration := a.points[len(a.points)-1].Time().Sub(a.points[0].Time())
	return NewSpeedFromDistanceAndDuration(distance, duration)
}

func (a *SpeedMovingAverage) distanceBetweenPoints() Distance {
	var distance Distance
	for i := 0; i < len(a.points)-1; i++ {
		distance = distance.Add(a.points[i].Position().Distance(a.points[i+1].Position()))
	}
	return distance
}
