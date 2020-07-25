package domain

import (
	"fmt"
	"sort"
	"time"

	"github.com/boreq/velo/internal/eventsourcing"
	"github.com/boreq/errors"
)

type Route struct {
	uuid   RouteUUID
	points []Point

	es eventsourcing.EventSourcing
}

func NewRoute(uuid RouteUUID, points []Point) (*Route, error) {
	if uuid.IsZero() {
		return nil, errors.New("zero value of uuid")
	}

	if len(points) == 0 {
		return nil, errors.New("missing points")
	}

	return &Route{
		uuid:   uuid,
		points: NormaliseRoutePoints(points),
	}, nil
}

func (r Route) UUID() RouteUUID {
	return r.uuid
}

func (r Route) Points() []Point {
	var points []Point

	for _, point := range r.points {
		points = append(points, point)
	}

	return points
}

func (r *Route) update(event eventsourcing.Event) error {
	switch e := event.(type) {
	case RouteCreated:
		r.handleRouteCreated(e)
	default:
		return fmt.Errorf("unknown event type '%T'", event)
	}

	return r.es.Record(event)
}

func (r *Route) handleRouteCreated(event RouteCreated) {
	r.uuid = event.UUID
	r.points = event.Points
}

type RouteCreated struct {
	UUID   RouteUUID
	Points []Point
}

func (r RouteCreated) EventType() eventsourcing.EventType {
	return "RouteCreated_v1"
}

// Points that occur more often than that will be dropped. Should probably be
// slightly lower than the desired 10 seconds just to account for tiny
// precision problems.
const intervalBetweenPoints = 9 * time.Second

func NormaliseRoutePoints(points []Point) []Point {
	sort.Slice(points, func(i, j int) bool {
		return points[i].Time().Before(points[j].Time())
	})

	var normalised []Point

	for _, point := range points {
		if len(normalised) == 0 {
			normalised = append(normalised, point)
		} else {
			previous := normalised[len(normalised)-1]
			if shouldAdd(previous, point) {
				normalised = append(normalised, point)
			}
		}
	}

	return normalised
}

func shouldAdd(previous Point, next Point) bool {
	if !next.Time().Before(previous.Time().Add(intervalBetweenPoints)) {
		return false
	}
	return true
}
