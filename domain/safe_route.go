package domain

import (
	"time"

	"github.com/boreq/errors"
)

const (
	notMovingSpeedThreshold = 0.25 // meters per second
	averageSpeedWindow      = 60 * time.Second
)

type SafeRoute struct {
	uuid        RouteUUID
	safePoints  []AnnotatedPoint
	distance    Distance
	timeMoving  time.Duration
	timeStarted time.Time
	timeEnded   time.Time
}

func NewSafeRoute(route *Route, privacyZones []*PrivacyZone) (*SafeRoute, error) {
	for _, privacyZone := range privacyZones {
		if privacyZone.IsZero() {
			return nil, errors.New("zero value of privacy zone")
		}
	}

	points := route.Points()

	var result []AnnotatedPoint

	avg := NewSpeedMovingAverage(averageSpeedWindow)
	var distance Distance
	var timeMoving time.Duration

	for i := range points {
		point := points[i]

		if i > 0 {
			previousPoint := points[i-1]
			distance = distance.Add(
				point.Position().Distance(previousPoint.Position()),
			)
		}

		avg.AddPoint(point)

		speed, err := avg.Speed()
		if err != nil {
			return nil, errors.Wrap(err, "could not retrieve the average speed")
		}

		annotatedPoint, err := NewAnnotatedPoint(point, speed, distance)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a new annotated point")
		}

		result = append(result, annotatedPoint)

		// detect that this is an auto pause
		if speed.Float64() < notMovingSpeedThreshold {
			if i > 0 {
				result[i-1].SetSpeedToZero()
				result[i].SetSpeedToZero()
			}
		}

		if i > 0 {
			a := result[i]
			b := result[i-1]
			if !(a.Speed().IsZero() && b.Speed().IsZero()) {
				timeMoving += a.Time().Sub(b.Time())
			}
		}
	}

	return &SafeRoute{
		uuid:        route.UUID(),
		distance:    distance,
		safePoints:  makePointsSafe(result, privacyZones),
		timeMoving:  timeMoving,
		timeStarted: points[0].Time(),
		timeEnded:   points[len(points)-1].Time(),
	}, nil
}

func (r SafeRoute) UUID() RouteUUID {
	return r.uuid
}

func (r SafeRoute) Points() []AnnotatedPoint {
	points := make([]AnnotatedPoint, len(r.safePoints))
	copy(points, r.safePoints)
	return points
}

func (r SafeRoute) TimeStarted() time.Time {
	return r.timeStarted
}

func (r SafeRoute) TimeEnded() time.Time {
	return r.timeEnded
}

func (r SafeRoute) TimeMoving() time.Duration {
	return r.timeMoving
}

func (r SafeRoute) Distance() Distance {
	return r.distance
}

func makePointsSafe(points []AnnotatedPoint, privacyZones []*PrivacyZone) []AnnotatedPoint {
	var safePoints []AnnotatedPoint

	for _, point := range points {
		if !positionIsWithinPrivacyZones(point.Position(), privacyZones) {
			safePoints = append(safePoints, point)
		}
	}

	return safePoints
}

func positionIsWithinPrivacyZones(position Position, privacyZones []*PrivacyZone) bool {
	for _, privacyZone := range privacyZones {
		if privacyZone.Circle().Contains(position) {
			return true
		}
	}
	return false
}

func distanceBetweenPoints(points []Point) Distance {
	var distance Distance
	for i := 0; i < len(points)-1; i++ {
		distance = distance.Add(points[i].Position().Distance(points[i+1].Position()))
	}
	return distance
}
