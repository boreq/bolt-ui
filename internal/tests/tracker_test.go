package tests

import (
	"sort"
	"testing"

	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/internal/fixture"
	"github.com/boreq/velo/internal/wire"
	"github.com/stretchr/testify/require"
)

func TestAddActivity(t *testing.T) {
	tr, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	gpxFile, cleanupFile := fixture.TestDataFile(t, "data/strava_export.gpx")
	defer cleanupFile()

	userUUID := auth.MustNewUserUUID("user-uuid")

	cmd := tracker.AddActivity{
		RouteFile: gpxFile,
		UserUUID:  userUUID,
	}

	activityUUID, err := tr.AddActivity.Execute(cmd)
	require.NoError(t, err)

	require.False(t, activityUUID.IsZero())

	result, err := tr.GetActivity.Execute(
		tracker.GetActivity{
			ActivityUUID: activityUUID,
		},
	)
	require.NoError(t, err)

	require.False(t, result.Activity.UUID().IsZero())
	require.False(t, result.Activity.RouteUUID().IsZero())
	require.False(t, result.Activity.UserUUID().IsZero())

	require.False(t, result.Route.UUID().IsZero())
	require.NotEmpty(t, result.Route.Points())

	activities, err := tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID: userUUID,
		},
	)
	require.NoError(t, err)
	require.Len(t, activities.Activities, 1)
	require.Equal(t, activities.Activities[0].UUID(), activityUUID)
	require.False(t, activities.HasNext)
	require.False(t, activities.HasPrev)
}

func TestListUserActivities(t *testing.T) {
	tr, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	userUUID := auth.MustNewUserUUID("user-uuid")

	var activityUUIDs []domain.ActivityUUID

	for i := 0; i < 30; i++ {
		gpxFile, cleanupFile := fixture.TestDataFile(t, "data/strava_export.gpx")
		defer cleanupFile()

		cmd := tracker.AddActivity{
			RouteFile: gpxFile,
			UserUUID:  userUUID,
		}

		activityUUID, err := tr.AddActivity.Execute(cmd)
		require.NoError(t, err)

		require.False(t, activityUUID.IsZero())

		activityUUIDs = append(activityUUIDs, activityUUID)
	}

	sort.Slice(activityUUIDs, func(i, j int) bool {
		return activityUUIDs[i].String() > activityUUIDs[j].String()
	})

	page1 := activityUUIDs[0:10]
	page2 := activityUUIDs[10:20]
	page3 := activityUUIDs[20:30]

	// page1 (initial)
	activities, err := tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID: userUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page1, toUUIDs(activities.Activities))
	require.False(t, activities.HasPrev)
	require.True(t, activities.HasNext)

	afterUUID := activities.Activities[len(activities.Activities)-1].UUID()

	// page2 (from page1)
	activities, err = tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID:   userUUID,
			StartAfter: &afterUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page2, toUUIDs(activities.Activities))
	require.True(t, activities.HasPrev)
	require.True(t, activities.HasNext)

	afterUUID = activities.Activities[len(activities.Activities)-1].UUID()

	// page3 (from page2)
	activities, err = tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID:   userUUID,
			StartAfter: &afterUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page3, toUUIDs(activities.Activities))
	require.True(t, activities.HasPrev)
	require.False(t, activities.HasNext)

	beforeUUID := activities.Activities[0].UUID()

	// page2 (from page3)
	activities, err = tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID:    userUUID,
			StartBefore: &beforeUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page2, toUUIDs(activities.Activities))
	require.True(t, activities.HasPrev)
	require.True(t, activities.HasNext)

	beforeUUID = activities.Activities[0].UUID()

	// page1 (from page2)
	activities, err = tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID:    userUUID,
			StartBefore: &beforeUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page1, toUUIDs(activities.Activities))
	require.False(t, activities.HasPrev)
	require.True(t, activities.HasNext)
}

func NewTracker(t *testing.T) (*tracker.Tracker, fixture.CleanupFunc) {
	db, cleanup := fixture.Bolt(t)

	tr, err := wire.BuildTrackerForTest(db)
	if err != nil {
		t.Fatal(err)
	}

	return tr, cleanup
}

func toUUIDs(activities []*domain.Activity) []domain.ActivityUUID {
	var uuids []domain.ActivityUUID

	for _, acitivity := range activities {
		uuids = append(uuids, acitivity.UUID())
	}

	return uuids
}
