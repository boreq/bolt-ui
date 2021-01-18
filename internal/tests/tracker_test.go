package tests

import (
	"errors"
	"sort"
	"testing"

	appAuth "github.com/boreq/velo/application/auth"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/internal/fixture"
	"github.com/boreq/velo/internal/wire"
	"github.com/stretchr/testify/require"
)

func TestAddActivity(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	gpxFile, cleanupFile := fixture.TestDataFile(t, "data/strava_export.gpx")
	defer cleanupFile()

	userUUID := auth.MustNewUserUUID("user-uuid")
	visibility := domain.PublicActivityVisibility
	title := domain.MustNewActivityTitle("title")

	testTracker.UserRepository.Users[userUUID] = appAuth.User{
		Username: "username",
	}

	cmd := tracker.AddActivity{
		RouteFile:  gpxFile,
		UserUUID:   userUUID,
		Visibility: visibility,
		Title:      title,
	}

	activityUUID, err := tr.AddActivity.Execute(cmd)
	require.NoError(t, err)

	require.False(t, activityUUID.IsZero())

	// test get
	result, err := tr.GetActivity.Execute(
		tracker.GetActivity{
			ActivityUUID: activityUUID,
		},
	)
	require.NoError(t, err)

	require.False(t, result.Activity.UUID().IsZero())
	require.False(t, result.Activity.RouteUUID().IsZero())
	require.Equal(t, userUUID, result.Activity.UserUUID())
	require.Equal(t, visibility, result.Activity.Visibility())
	require.Equal(t, title, result.Activity.Title())

	require.False(t, result.Route.UUID().IsZero())
	require.NotEmpty(t, result.Route.Points())

	require.NotEmpty(t, result.User.Username)

	// test list
	activities, err := tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID: userUUID,
		},
	)
	require.NoError(t, err)
	require.Len(t, activities.Activities, 1)
	require.Equal(t, activities.Activities[0].Activity.UUID(), activityUUID)
	require.NotEmpty(t, activities.Activities[0].User.Username)
	require.False(t, activities.HasNext)
	require.False(t, activities.HasPrevious)
}

func TestListUserActivities(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	userUUID := auth.MustNewUserUUID("user-uuid")

	testTracker.UserRepository.Users[userUUID] = appAuth.User{
		Username: "username",
	}

	var activityUUIDs []domain.ActivityUUID

	for i := 0; i < 30; i++ {
		gpxFile, cleanupFile := fixture.TestDataFile(t, "data/strava_export.gpx")
		defer cleanupFile()

		cmd := tracker.AddActivity{
			RouteFile:  gpxFile,
			UserUUID:   userUUID,
			Visibility: domain.PublicActivityVisibility,
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
	require.False(t, activities.HasPrevious)
	require.True(t, activities.HasNext)

	afterUUID := activities.Activities[len(activities.Activities)-1].Activity.UUID()

	// page2 (from page1)
	activities, err = tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID:   userUUID,
			StartAfter: &afterUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page2, toUUIDs(activities.Activities))
	require.True(t, activities.HasPrevious)
	require.True(t, activities.HasNext)

	afterUUID = activities.Activities[len(activities.Activities)-1].Activity.UUID()

	// page3 (from page2)
	activities, err = tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID:   userUUID,
			StartAfter: &afterUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page3, toUUIDs(activities.Activities))
	require.True(t, activities.HasPrevious)
	require.False(t, activities.HasNext)

	beforeUUID := activities.Activities[0].Activity.UUID()

	// page2 (from page3)
	activities, err = tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID:    userUUID,
			StartBefore: &beforeUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page2, toUUIDs(activities.Activities))
	require.True(t, activities.HasPrevious)
	require.True(t, activities.HasNext)

	beforeUUID = activities.Activities[0].Activity.UUID()

	// page1 (from page2)
	activities, err = tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID:    userUUID,
			StartBefore: &beforeUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page1, toUUIDs(activities.Activities))
	require.False(t, activities.HasPrevious)
	require.True(t, activities.HasNext)
}

func TestPermissions(t *testing.T) {
	testCases := []struct {
		Visibility domain.ActivityVisibility

		UnauthorisedCanView bool
		OtherCanView        bool
		OwnerCanView        bool

		UnauthorisedCanList bool
		OtherCanList        bool
		OwnerCanList        bool
	}{
		{
			Visibility: domain.PublicActivityVisibility,

			UnauthorisedCanView: true,
			OtherCanView:        true,
			OwnerCanView:        true,

			UnauthorisedCanList: true,
			OtherCanList:        true,
			OwnerCanList:        true,
		},
		{
			Visibility: domain.UnlistedActivityVisibility,

			UnauthorisedCanView: true,
			OtherCanView:        true,
			OwnerCanView:        true,

			UnauthorisedCanList: false,
			OtherCanList:        false,
			OwnerCanList:        true,
		},
		{
			Visibility: domain.PrivateActivityVisibility,

			UnauthorisedCanView: false,
			OtherCanView:        false,
			OwnerCanView:        true,

			UnauthorisedCanList: false,
			OtherCanList:        false,
			OwnerCanList:        true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Visibility.String(), func(t *testing.T) {
			testTracker, cleanupTracker := NewTracker(t)
			defer cleanupTracker()

			tr := testTracker.Tracker

			gpxFile, cleanupFile := fixture.TestDataFile(t, "data/strava_export.gpx")
			defer cleanupFile()

			user := auth.ReadUser{
				UUID: auth.MustNewUserUUID("user-uuid"),
			}

			otherUser := auth.ReadUser{
				UUID: auth.MustNewUserUUID("other-user-uuid"),
			}

			testTracker.UserRepository.Users[user.UUID] = appAuth.User{
				Username: "username",
			}

			cmd := tracker.AddActivity{
				RouteFile:  gpxFile,
				UserUUID:   user.UUID,
				Visibility: testCase.Visibility,
			}

			activityUUID, err := tr.AddActivity.Execute(cmd)
			require.NoError(t, err)

			t.Run("get unauthorised", func(t *testing.T) {
				_, err = tr.GetActivity.Execute(
					tracker.GetActivity{
						ActivityUUID: activityUUID,
						AsUser:       nil,
					},
				)

				if testCase.UnauthorisedCanView {
					require.NoError(t, err)
				} else {
					require.True(t, errors.Is(err, tracker.ErrGettingActivityForbidden))
				}
			})

			t.Run("get other", func(t *testing.T) {
				_, err = tr.GetActivity.Execute(
					tracker.GetActivity{
						ActivityUUID: activityUUID,
						AsUser:       &otherUser,
					},
				)

				if testCase.OtherCanView {
					require.NoError(t, err)
				} else {
					require.True(t, errors.Is(err, tracker.ErrGettingActivityForbidden))
				}
			})

			t.Run("get owner", func(t *testing.T) {
				_, err = tr.GetActivity.Execute(
					tracker.GetActivity{
						ActivityUUID: activityUUID,
						AsUser:       &user,
					},
				)

				if testCase.OwnerCanView {
					require.NoError(t, err)
				} else {
					require.True(t, errors.Is(err, tracker.ErrGettingActivityForbidden))
				}
			})

			t.Run("list unauthorised", func(t *testing.T) {
				result, err := tr.ListUserActivities.Execute(
					tracker.ListUserActivities{
						UserUUID: user.UUID,
						AsUser:   nil,
					},
				)
				require.NoError(t, err)

				if testCase.UnauthorisedCanList {
					require.NotEmpty(t, result.Activities)
				} else {
					require.Empty(t, result.Activities)
				}
			})

			t.Run("list other", func(t *testing.T) {
				result, err := tr.ListUserActivities.Execute(
					tracker.ListUserActivities{
						UserUUID: user.UUID,
						AsUser:   &otherUser,
					},
				)
				require.NoError(t, err)

				if testCase.OtherCanList {
					require.NotEmpty(t, result.Activities)
				} else {
					require.Empty(t, result.Activities)
				}
			})

			t.Run("list owner", func(t *testing.T) {
				result, err := tr.ListUserActivities.Execute(
					tracker.ListUserActivities{
						UserUUID: user.UUID,
						AsUser:   &user,
					},
				)
				require.NoError(t, err)

				if testCase.OwnerCanList {
					require.NotEmpty(t, result.Activities)
				} else {
					require.Empty(t, result.Activities)
				}
			})
		})
	}

}

func NewTracker(t *testing.T) (wire.TestTracker, fixture.CleanupFunc) {
	db, cleanup := fixture.Bolt(t)

	tr, err := wire.BuildTrackerForTest(db)
	if err != nil {
		t.Fatal(err)
	}

	return tr, cleanup
}

func toUUIDs(activities []tracker.Activity) []domain.ActivityUUID {
	var uuids []domain.ActivityUUID

	for _, acitivity := range activities {
		uuids = append(uuids, acitivity.Activity.UUID())
	}

	return uuids
}
