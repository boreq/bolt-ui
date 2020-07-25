package domain_test

import (
	"testing"
	"time"

	"github.com/boreq/eggplant/domain"
	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {
	uuid := domain.MustNewRouteUUID("route-uuid")
	points := somePoints()

	route, err := domain.NewRoute(uuid, points)
	require.NoError(t, err)

	require.Equal(t, uuid, route.UUID())
	require.NotZero(t, len(route.Points()))
}

func somePoints() []domain.Point {
	date := time.Date(1954, time.June, 7, 12, 0, 0, 0, time.UTC)

	var points []domain.Point
	for i := 0; i < 10; i++ {
		p := domain.MustNewPoint(
			date.Add(time.Duration(i)*time.Minute),
			somePosition(),
			someAltitude(),
		)
		points = append(points, p)
	}
	return points
}

func someAltitude() domain.Altitude {
	return domain.NewAltitude(10)
}

func somePosition() domain.Position {
	return domain.NewPosition(
		domain.MustNewLatitude(10),
		domain.MustNewLongitude(10),
	)
}
