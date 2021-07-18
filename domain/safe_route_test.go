package domain_test

import (
	"testing"
	"time"

	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/stretchr/testify/require"
)

func TestSafeRoute(t *testing.T) {
	p1, err := domain.NewPoint(
		time.Date(2021, 01, 16, 11, 12, 13, 0, time.UTC),
		domain.NewPosition(
			domain.MustNewLatitude(10),
			domain.MustNewLongitude(10),
		),
		domain.MustNewAltitude(10),
	)
	require.NoError(t, err)

	p2, err := domain.NewPoint(
		time.Date(2021, 01, 16, 11, 13, 13, 0, time.UTC),
		domain.NewPosition(
			domain.MustNewLatitude(20),
			domain.MustNewLongitude(20),
		),
		domain.MustNewAltitude(10),
	)
	require.NoError(t, err)

	p3, err := domain.NewPoint(
		time.Date(2021, 01, 16, 11, 14, 13, 0, time.UTC),
		domain.NewPosition(
			domain.MustNewLatitude(30),
			domain.MustNewLongitude(30),
		),
		domain.MustNewAltitude(10),
	)
	require.NoError(t, err)

	route, err := domain.NewRoute(
		domain.MustNewRouteUUID("route-uuid"),
		[]domain.Point{p1, p2, p3},
	)
	require.NoError(t, err)

	require.Len(t, route.Points(), 3)

	privacyZone, err := domain.NewPrivacyZone(
		domain.MustNewPrivacyZoneUUID("privacy-zone-uuid"),
		auth.MustNewUserUUID("user-uuid"),
		domain.NewPosition(
			domain.MustNewLatitude(20),
			domain.MustNewLongitude(20),
		),
		domain.MustNewCircle(
			domain.NewPosition(
				domain.MustNewLatitude(20),
				domain.MustNewLongitude(20),
			),
			1000,
		),
		domain.MustNewPrivacyZoneName("some-privacy-zone"),
	)
	require.NoError(t, err)

	safeRoute, err := domain.NewSafeRoute(route, []*domain.PrivacyZone{privacyZone})
	require.NoError(t, err)

	require.Less(t, len(safeRoute.Points()), len(route.Points()))
}

func TestSafeRouteDistance(t *testing.T) {
	p1, err := domain.NewPoint(
		time.Date(2021, 01, 16, 11, 12, 13, 0, time.UTC),
		domain.NewPosition(
			domain.MustNewLatitude(10),
			domain.MustNewLongitude(10),
		),
		domain.MustNewAltitude(10),
	)
	require.NoError(t, err)

	p2, err := domain.NewPoint(
		time.Date(2021, 01, 16, 11, 13, 13, 0, time.UTC),
		domain.NewPosition(
			domain.MustNewLatitude(20),
			domain.MustNewLongitude(20),
		),
		domain.MustNewAltitude(10),
	)
	require.NoError(t, err)

	p3, err := domain.NewPoint(
		time.Date(2021, 01, 16, 11, 14, 13, 0, time.UTC),
		domain.NewPosition(
			domain.MustNewLatitude(30),
			domain.MustNewLongitude(30),
		),
		domain.MustNewAltitude(10),
	)
	require.NoError(t, err)

	route, err := domain.NewRoute(
		domain.MustNewRouteUUID("route-uuid"),
		[]domain.Point{p1, p2, p3},
	)
	require.NoError(t, err)

	require.Len(t, route.Points(), 3)

	safeRoute, err := domain.NewSafeRoute(route, nil)
	require.NoError(t, err)

	expectedDistance := p1.Position().Distance(p2.Position()).Add(p2.Position().Distance(p3.Position()))
	require.InEpsilon(t, expectedDistance.Float64(), safeRoute.Distance().Float64(), 0.0001)
}
