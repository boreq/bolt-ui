package domain_test

import (
	"testing"

	"github.com/boreq/velo/domain"
	"github.com/stretchr/testify/require"
)

func TestDistance(t *testing.T) {
	p1 := domain.NewPosition(
		domain.MustNewLatitude(50.061389),
		domain.MustNewLongitude(19.937222),
	)

	p2 := domain.NewPosition(
		domain.MustNewLatitude(52.52),
		domain.MustNewLongitude(13.405),
	)

	d1 := p1.Distance(p2)
	require.NotZero(t, d1)

	d2 := p2.Distance(p1)
	require.NotZero(t, d2)

	require.Equal(t, d1, d2)
}
