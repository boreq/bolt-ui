package domain_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/boreq/velo/domain"
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

func TestNormaliseRoutePoints(t *testing.T) {
	date := time.Date(1954, time.June, 7, 12, 0, 0, 0, time.UTC)

	p1 := somePoint(date)
	p2 := somePoint(date.Add(2 * time.Second))
	p3 := somePoint(date.Add(5 * time.Second))

	testCases := []struct {
		Input  []domain.Point
		Output []domain.Point
	}{
		{
			Input:  nil,
			Output: nil,
		},
		{
			Input: []domain.Point{
				p1,
			},
			Output: []domain.Point{
				p1,
			},
		},
		{
			Input: []domain.Point{
				p1,
				p2,
			},
			Output: []domain.Point{
				p1,
				p2,
			},
		},
		{
			Input: []domain.Point{
				p1,
				p2,
				p3,
			},
			Output: []domain.Point{
				p1,
				p3,
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			output := domain.NormaliseRoutePoints(testCase.Input)
			require.Equal(t, testCase.Output, output)
		})
	}
}

func somePoints() []domain.Point {
	date := time.Date(1954, time.June, 7, 12, 0, 0, 0, time.UTC)

	var points []domain.Point
	for i := 0; i < 10; i++ {
		t := date.Add(time.Duration(i) * time.Minute)
		points = append(points, somePoint(t))
	}
	return points
}

func somePoint(t time.Time) domain.Point {
	return domain.MustNewPoint(
		t,
		somePosition(),
		someAltitude(),
	)
}

func someAltitude() domain.Altitude {
	return domain.MustNewAltitude(10)
}

func somePosition() domain.Position {
	return domain.NewPosition(
		domain.MustNewLatitude(10),
		domain.MustNewLongitude(10),
	)
}
