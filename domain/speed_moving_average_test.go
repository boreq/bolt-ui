package domain_test

import (
	"testing"
	"time"

	"github.com/boreq/velo/domain"
	"github.com/stretchr/testify/require"
)

func TestNewSpeedMovingAverage(t *testing.T) {
	date := someTime()

	avg := domain.NewSpeedMovingAverage(5 * time.Second)

	spd, err := avg.Speed()
	require.NoError(t, err)
	require.Zero(t, spd.Float64())

	p1 := domain.MustNewPoint(
		date,
		domain.NewPosition(
			domain.MustNewLatitude(10),
			domain.MustNewLongitude(10),
		),
		someAltitude(),
	)

	avg.AddPoint(p1)

	spd, err = avg.Speed()
	require.NoError(t, err)
	require.Zero(t, spd.Float64())

	p2 := domain.MustNewPoint(
		date.Add(1*time.Second),
		domain.NewPosition(
			domain.MustNewLatitude(20),
			domain.MustNewLongitude(20),
		),
		someAltitude(),
	)

	avg.AddPoint(p2)

	spd, err = avg.Speed()
	require.NoError(t, err)
	require.NotZero(t, spd.Float64())
}

func TestNewSpeedMovingAverage_points_further_apart_than_duration(t *testing.T) {
	date := someTime()

	avg := domain.NewSpeedMovingAverage(5 * time.Second)

	spd, err := avg.Speed()
	require.NoError(t, err)
	require.Zero(t, spd.Float64())

	p1 := domain.MustNewPoint(
		date,
		domain.NewPosition(
			domain.MustNewLatitude(10),
			domain.MustNewLongitude(10),
		),
		someAltitude(),
	)

	avg.AddPoint(p1)

	spd, err = avg.Speed()
	require.NoError(t, err)
	require.Zero(t, spd.Float64())

	p2 := domain.MustNewPoint(
		date.Add(10*time.Second),
		domain.NewPosition(
			domain.MustNewLatitude(20),
			domain.MustNewLongitude(20),
		),
		someAltitude(),
	)

	avg.AddPoint(p2)

	spd, err = avg.Speed()
	require.NoError(t, err)
	require.NotZero(t, spd.Float64())
}

func someTime() time.Time {
	return time.Date(1954, time.June, 7, 12, 0, 0, 0, time.UTC)
}
