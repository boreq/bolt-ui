package domain

import (
	"time"

	"github.com/boreq/errors"
)

type SafeRoute struct {
	uuid       RouteUUID
	points     []Point
	safePoints []Point
}

func NewSafeRoute(route *Route, privacyZones []*PrivacyZone) (*SafeRoute, error) {
	for _, privacyZone := range privacyZones {
		if privacyZone.IsZero() {
			return nil, errors.New("zero value of privacy zone")
		}
	}

	return &SafeRoute{
		uuid:       route.UUID(),
		points:     route.Points(),
		safePoints: makePointsSafe(route.Points(), privacyZones),
	}, nil
}

func (r SafeRoute) UUID() RouteUUID {
	return r.uuid
}

func (r SafeRoute) Points() []Point {
	points := make([]Point, len(r.safePoints))
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

func (r SafeRoute) Distance() float64 {
	var distance float64
	for i := 0; i < len(r.points)-1; i++ {
		distance += r.points[i].Position().Distance(r.points[i+1].Position())
	}
	return distance
}

func makePointsSafe(points []Point, privacyZones []*PrivacyZone) []Point {
	var safePoints []Point

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
