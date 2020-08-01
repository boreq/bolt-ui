package tests

import (
	"testing"

	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/internal/fixture"
	"github.com/boreq/velo/internal/wire"
	"github.com/stretchr/testify/require"
)

func TestAddActivity(t *testing.T) {
	tr, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	gpxFile, cleanupFile := fixture.TestDataFile(t, "data/strava_export.gpx")
	defer cleanupFile()

	userUUID := domain.MustNewUserUUID("user-uuid")

	cmd := tracker.AddActivity{
		RouteFile: gpxFile,
		UserUUID:  userUUID,
	}

	activityUUID, err := tr.AddActivity.Execute(cmd)
	require.NoError(t, err)

	require.False(t, activityUUID.IsZero())

	result, err := tr.GetActivity.Execute(tracker.GetActivity{
		ActivityUUID: activityUUID,
	})
	require.NoError(t, err)

	require.False(t, result.Activity.UUID().IsZero())
	require.False(t, result.Activity.RouteUUID().IsZero())
	require.False(t, result.Activity.UserUUID().IsZero())

	require.False(t, result.Route.UUID().IsZero())
	require.NotEmpty(t, result.Route.Points())

	// todo check that the route can be listed in the user profile
}

func NewTracker(t *testing.T) (*tracker.Tracker, fixture.CleanupFunc) {
	db, cleanup := fixture.Bolt(t)

	tr, err := wire.BuildTrackerForTest(db)
	if err != nil {
		t.Fatal(err)
	}

	return tr, cleanup
}
