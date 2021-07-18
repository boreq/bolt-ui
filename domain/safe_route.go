package domain

import (
	"time"

	"github.com/boreq/errors"
)

type SafeRoute struct {
	uuid       RouteUUID
	points     []AnnotatedPoint
	safePoints []AnnotatedPoint
}

func NewSafeRoute(route *Route, privacyZones []*PrivacyZone) (*SafeRoute, error) {
	for _, privacyZone := range privacyZones {
		if privacyZone.IsZero() {
			return nil, errors.New("zero value of privacy zone")
		}
	}

	annotatedPoints, err := toAnnotatedPoints(route.Points())
	if err != nil {
		return nil, errors.Wrap(err, "could not convert points to annotated points")
	}

	return &SafeRoute{
		uuid:       route.UUID(),
		points:     annotatedPoints,
		safePoints: makePointsSafe(annotatedPoints, privacyZones),
	}, nil
}

func (r SafeRoute) UUID() RouteUUID {
	return r.uuid
}

func (r SafeRoute) Points() []Point {
	var points []Point
	for _, point := range r.safePoints {
		points = append(points, point.Point())
	}
	return points
}

func (r SafeRoute) AnnotatedPoints() []AnnotatedPoint {
	points := make([]AnnotatedPoint, len(r.safePoints))
	copy(points, r.safePoints)
	return points
}

func (r SafeRoute) TimeStarted() time.Time {
	return r.points[0].Time()
}

func (r SafeRoute) TimeEnded() time.Time {
	return r.points[len(r.points)-1].Time()
}

func (r SafeRoute) TimeMoving() time.Duration {
	return r.TimeEnded().Sub(r.TimeStarted())
}

func (r SafeRoute) Distance() Distance {
	return distanceBetweenPoints(r.Points())
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

const notMovingSpeedThreshold = 0.25 // meters per second

func toAnnotatedPoints(points []Point) ([]AnnotatedPoint, error) {
	var result []AnnotatedPoint

	avg := NewSpeedMovingAverage(60 * time.Second)
	var cumulativeDistance Distance

	for i := range points {
		point := points[i]

		if i > 0 {
			previousPoint := points[i-1]
			cumulativeDistance = cumulativeDistance.Add(
				point.Position().Distance(previousPoint.Position()),
			)
		}

		avg.AddPoint(point)

		speed, err := avg.Speed()
		if err != nil {
			return nil, errors.Wrap(err, "could not retrieve the average speed")
		}

		annotatedPoint, err := NewAnnotatedPoint(point, speed, cumulativeDistance)
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
	}

	return result, nil
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
