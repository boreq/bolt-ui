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

const testRouteFile = "data/strava_export_route.gpx"
const testStravaExportFile = "data/strava_export.zip"
const testStravaExportFileWithFit = "data/strava_export_with_fit.zip"

func TestAddActivity(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	gpxFile, cleanupFile := fixture.TestDataFile(t, testRouteFile)
	defer cleanupFile()

	visibility := domain.PublicActivityVisibility
	title := domain.MustNewActivityTitle("title")

	user := mockUser()
	testTracker.UserRepository.Users[user.UUID()] = user

	cmd := tracker.AddActivity{
		RouteFile:       gpxFile,
		RouteFileFormat: tracker.RouteFileFormatGpx,
		UserUUID:        user.UUID(),
		Visibility:      visibility,
		Title:           title,
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
	require.Equal(t, user.UUID(), result.Activity.UserUUID())
	require.Equal(t, visibility, result.Activity.Visibility())
	require.Equal(t, title, result.Activity.Title())

	require.False(t, result.Route.UUID().IsZero())
	require.NotEmpty(t, result.Route.Points())

	require.NotEmpty(t, result.User.Username)

	// test list
	activities, err := tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID: user.UUID(),
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

	user := mockUser()
	testTracker.UserRepository.Users[user.UUID()] = user

	var activityUUIDs []domain.ActivityUUID

	for i := 0; i < 30; i++ {
		gpxFile, cleanupFile := fixture.TestDataFile(t, testRouteFile)
		defer cleanupFile()

		cmd := tracker.AddActivity{
			RouteFile:       gpxFile,
			RouteFileFormat: tracker.RouteFileFormatGpx,
			UserUUID:        user.UUID(),
			Visibility:      domain.PublicActivityVisibility,
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
			UserUUID: user.UUID(),
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
			UserUUID:   user.UUID(),
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
			UserUUID:   user.UUID(),
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
			UserUUID:    user.UUID(),
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
			UserUUID:    user.UUID(),
			StartBefore: &beforeUUID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, page1, toUUIDs(activities.Activities))
	require.False(t, activities.HasPrevious)
	require.True(t, activities.HasNext)
}

func TestActivityPermissions(t *testing.T) {
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

			gpxFile, cleanupFile := fixture.TestDataFile(t, testRouteFile)
			defer cleanupFile()

			user := mockUser()
			testTracker.UserRepository.Users[user.UUID()] = user

			readUser := auth.ReadUser{
				UUID: user.UUID(),
			}

			readOtherUser := auth.ReadUser{
				UUID: auth.MustNewUserUUID("other-user-uuid"),
			}

			cmd := tracker.AddActivity{
				RouteFile:       gpxFile,
				RouteFileFormat: tracker.RouteFileFormatGpx,
				UserUUID:        user.UUID(),
				Visibility:      testCase.Visibility,
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
					require.ErrorIs(t, err, tracker.ErrGettingActivityForbidden)
				}
			})

			t.Run("get other", func(t *testing.T) {
				_, err = tr.GetActivity.Execute(
					tracker.GetActivity{
						ActivityUUID: activityUUID,
						AsUser:       &readOtherUser,
					},
				)

				if testCase.OtherCanView {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, tracker.ErrGettingActivityForbidden)
				}
			})

			t.Run("get owner", func(t *testing.T) {
				_, err = tr.GetActivity.Execute(
					tracker.GetActivity{
						ActivityUUID: activityUUID,
						AsUser:       &readUser,
					},
				)

				if testCase.OwnerCanView {
					require.NoError(t, err)
				} else {
					require.ErrorIs(t, err, tracker.ErrGettingActivityForbidden)
				}
			})

			t.Run("list unauthorised", func(t *testing.T) {
				result, err := tr.ListUserActivities.Execute(
					tracker.ListUserActivities{
						UserUUID: user.UUID(),
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
						UserUUID: user.UUID(),
						AsUser:   &readOtherUser,
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
						UserUUID: user.UUID(),
						AsUser:   &readUser,
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

func TestEditActivityPermissions(t *testing.T) {
	user := auth.ReadUser{
		UUID: auth.MustNewUserUUID("user-uuid"),
	}

	otherUser := auth.ReadUser{
		UUID: auth.MustNewUserUUID("other-user-uuid"),
	}

	testCases := []struct {
		Name    string
		User    *auth.ReadUser
		CanEdit bool
	}{
		{
			Name:    "unauthorized_user",
			User:    nil,
			CanEdit: false,
		},
		{
			Name:    "other_user",
			User:    &otherUser,
			CanEdit: false,
		},
		{
			Name:    "user",
			User:    &user,
			CanEdit: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testTracker, cleanupTracker := NewTracker(t)
			defer cleanupTracker()

			tr := testTracker.Tracker

			gpxFile, cleanupFile := fixture.TestDataFile(t, testRouteFile)
			defer cleanupFile()

			user := mockUser()
			testTracker.UserRepository.Users[user.UUID()] = user

			cmd := tracker.AddActivity{
				RouteFile:       gpxFile,
				RouteFileFormat: tracker.RouteFileFormatGpx,
				UserUUID:        user.UUID(),
				Title:           domain.MustNewActivityTitle("title"),
				Visibility:      domain.PublicActivityVisibility,
			}

			activityUUID, err := tr.AddActivity.Execute(cmd)
			require.NoError(t, err)

			err = tr.EditActivity.Execute(
				tracker.EditActivity{
					ActivityUUID: activityUUID,
					AsUser:       testCase.User,
					Title:        domain.MustNewActivityTitle("new-title"),
					Visibility:   domain.PrivateActivityVisibility,
				},
			)

			if testCase.CanEdit {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tracker.ErrEditingActivityForbidden)
			}
		})
	}
}

func TestEditActivity(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	gpxFile, cleanupFile := fixture.TestDataFile(t, testRouteFile)
	defer cleanupFile()

	user := mockUser()
	testTracker.UserRepository.Users[user.UUID()] = user

	readUser := user.AsReadUser()

	initialTitle := domain.MustNewActivityTitle("title")
	initialVisibility := domain.PublicActivityVisibility

	cmd := tracker.AddActivity{
		RouteFile:       gpxFile,
		RouteFileFormat: tracker.RouteFileFormatGpx,
		UserUUID:        user.UUID(),
		Title:           initialTitle,
		Visibility:      initialVisibility,
	}

	activityUUID, err := tr.AddActivity.Execute(cmd)
	require.NoError(t, err)

	activity, err := tr.GetActivity.Execute(
		tracker.GetActivity{
			ActivityUUID: activityUUID,
			AsUser:       &readUser,
		},
	)
	require.NoError(t, err)

	require.Equal(t, initialTitle, activity.Activity.Title())
	require.Equal(t, initialVisibility, activity.Activity.Visibility())

	newTitle := domain.MustNewActivityTitle("new-title")
	newVisibility := domain.PrivateActivityVisibility

	require.NotEqual(t, initialTitle, newTitle)
	require.NotEqual(t, initialVisibility, newVisibility)

	err = tr.EditActivity.Execute(
		tracker.EditActivity{
			ActivityUUID: activityUUID,
			AsUser:       &readUser,
			Title:        newTitle,
			Visibility:   newVisibility,
		},
	)
	require.NoError(t, err)

	activity, err = tr.GetActivity.Execute(
		tracker.GetActivity{
			ActivityUUID: activityUUID,
			AsUser:       &readUser,
		},
	)
	require.NoError(t, err)

	require.Equal(t, newTitle, activity.Activity.Title())
	require.Equal(t, newVisibility, activity.Activity.Visibility())
}

func TestEditActivityWithoutChangesShouldNotReturnAnError(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	gpxFile, cleanupFile := fixture.TestDataFile(t, testRouteFile)
	defer cleanupFile()

	user := mockUser()
	testTracker.UserRepository.Users[user.UUID()] = user

	readUser := user.AsReadUser()

	initialTitle := domain.MustNewActivityTitle("title")
	initialVisibility := domain.PublicActivityVisibility

	cmd := tracker.AddActivity{
		RouteFile:       gpxFile,
		RouteFileFormat: tracker.RouteFileFormatGpx,
		UserUUID:        user.UUID(),
		Title:           initialTitle,
		Visibility:      initialVisibility,
	}

	activityUUID, err := tr.AddActivity.Execute(cmd)
	require.NoError(t, err)

	err = tr.EditActivity.Execute(
		tracker.EditActivity{
			ActivityUUID: activityUUID,
			AsUser:       &readUser,
			Title:        initialTitle,
			Visibility:   initialVisibility,
		},
	)
	require.NoError(t, err)
}

func TestDeleteActivityPermissions(t *testing.T) {
	user := auth.ReadUser{
		UUID: auth.MustNewUserUUID("user-uuid"),
	}

	otherUser := auth.ReadUser{
		UUID: auth.MustNewUserUUID("other-user-uuid"),
	}

	testCases := []struct {
		Name      string
		User      *auth.ReadUser
		CanDelete bool
	}{
		{
			Name:      "unauthorized_user",
			User:      nil,
			CanDelete: false,
		},
		{
			Name:      "other_user",
			User:      &otherUser,
			CanDelete: false,
		},
		{
			Name:      "user",
			User:      &user,
			CanDelete: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testTracker, cleanupTracker := NewTracker(t)
			defer cleanupTracker()

			tr := testTracker.Tracker

			gpxFile, cleanupFile := fixture.TestDataFile(t, testRouteFile)
			defer cleanupFile()

			user := mockUser()
			testTracker.UserRepository.Users[user.UUID()] = user

			cmd := tracker.AddActivity{
				RouteFile:       gpxFile,
				RouteFileFormat: tracker.RouteFileFormatGpx,
				UserUUID:        user.UUID(),
				Title:           domain.MustNewActivityTitle("title"),
				Visibility:      domain.PublicActivityVisibility,
			}

			activityUUID, err := tr.AddActivity.Execute(cmd)
			require.NoError(t, err)

			err = tr.DeleteActivity.Execute(
				tracker.DeleteActivity{
					ActivityUUID: activityUUID,
					AsUser:       testCase.User,
				},
			)

			if testCase.CanDelete {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tracker.ErrDeletingActivityForbidden)
			}
		})
	}
}

func TestDeleteActivity(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	gpxFile, cleanupFile := fixture.TestDataFile(t, testRouteFile)
	defer cleanupFile()

	user := mockUser()
	testTracker.UserRepository.Users[user.UUID()] = user
	readUser := user.AsReadUser()

	cmd := tracker.AddActivity{
		RouteFile:       gpxFile,
		RouteFileFormat: tracker.RouteFileFormatGpx,
		UserUUID:        user.UUID(),
		Title:           domain.MustNewActivityTitle("title"),
		Visibility:      domain.PublicActivityVisibility,
	}

	activityUUID, err := tr.AddActivity.Execute(cmd)
	require.NoError(t, err)

	_, err = tr.GetActivity.Execute(
		tracker.GetActivity{
			ActivityUUID: activityUUID,
			AsUser:       &readUser,
		},
	)
	require.NoError(t, err)

	activities, err := tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID: user.UUID(),
		},
	)
	require.NoError(t, err)
	require.Len(t, activities.Activities, 1)

	err = tr.DeleteActivity.Execute(
		tracker.DeleteActivity{
			ActivityUUID: activityUUID,
			AsUser:       &readUser,
		},
	)
	require.NoError(t, err)

	_, err = tr.GetActivity.Execute(
		tracker.GetActivity{
			ActivityUUID: activityUUID,
			AsUser:       &readUser,
		},
	)
	require.ErrorIs(t, err, tracker.ErrActivityNotFound)

	activities, err = tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID: user.UUID(),
		},
	)
	require.NoError(t, err)
	require.Len(t, activities.Activities, 0)
}

func TestDeleteActivityThatDoesNotExist(t *testing.T) {
	user := auth.ReadUser{
		UUID: auth.MustNewUserUUID("user-uuid"),
	}

	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	err := tr.DeleteActivity.Execute(
		tracker.DeleteActivity{
			ActivityUUID: domain.MustNewActivityUUID("activity-uuid"),
			AsUser:       &user,
		},
	)
	require.ErrorIs(t, err, tracker.ErrActivityNotFound)
}

func TestApplyPrivacyZones(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	gpxFile, cleanupFile := fixture.TestDataFile(t, testRouteFile)
	defer cleanupFile()

	user := mockUser()
	testTracker.UserRepository.Users[user.UUID()] = user
	readUser := user.AsReadUser()

	otherUser := auth.ReadUser{
		UUID: auth.MustNewUserUUID("other-user-uuid"),
	}

	cmd := tracker.AddActivity{
		RouteFile:       gpxFile,
		RouteFileFormat: tracker.RouteFileFormatGpx,
		UserUUID:        user.UUID(),
		Visibility:      domain.PublicActivityVisibility,
	}

	activityUUID, err := tr.AddActivity.Execute(cmd)
	require.NoError(t, err)

	position := domain.NewPosition(
		domain.MustNewLatitude(50.07357803907662),
		domain.MustNewLongitude(19.993221096609236),
	)

	circle := domain.MustNewCircle(
		domain.NewPosition(
			domain.MustNewLatitude(50.07357803907662),
			domain.MustNewLongitude(19.993221096609236),
		),
		500,
	)

	name := domain.MustNewPrivacyZoneName("Privacy zone")

	privacyZoneCmd := tracker.AddPrivacyZone{
		UserUUID: user.UUID(),
		Position: position,
		Circle:   circle,
		Name:     name,
	}

	_, err = tr.AddPrivacyZone.Execute(privacyZoneCmd)
	require.NoError(t, err)

	testCases := []struct {
		Name           string
		AsUser         *auth.ReadUser
		ExpectedPoints int
	}{
		{
			Name:           "unauthorized",
			AsUser:         nil,
			ExpectedPoints: 433,
		},
		{
			Name:           "other",
			AsUser:         &otherUser,
			ExpectedPoints: 433,
		},
		{
			Name:           "owner",
			AsUser:         &readUser,
			ExpectedPoints: 481,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Run("get", func(t *testing.T) {
				result, err := tr.GetActivity.Execute(
					tracker.GetActivity{
						ActivityUUID: activityUUID,
						AsUser:       testCase.AsUser,
					},
				)
				require.NoError(t, err)
				require.Len(t, result.Route.Points(), testCase.ExpectedPoints)
			})

			t.Run("list", func(t *testing.T) {
				result, err := tr.ListUserActivities.Execute(
					tracker.ListUserActivities{
						UserUUID: user.UUID(),
						AsUser:   testCase.AsUser,
					},
				)
				require.NoError(t, err)

				require.Len(t, result.Activities, 1)
				require.Len(t, result.Activities[0].Route.Points(), testCase.ExpectedPoints)
			})
		})
	}
}

func TestAddAndDeletePrivacyZone(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	user := mockUser()

	position := domain.NewPosition(
		domain.MustNewLatitude(48.20952),
		domain.MustNewLongitude(16.35618),
	)

	circle := domain.MustNewCircle(
		domain.NewPosition(
			domain.MustNewLatitude(48.21019),
			domain.MustNewLongitude(16.36163),
		),
		500,
	)

	name := domain.MustNewPrivacyZoneName("Privacy zone")

	cmd := tracker.AddPrivacyZone{
		UserUUID: user.UUID(),
		Position: position,
		Circle:   circle,
		Name:     name,
	}

	privacyZoneUUID, err := tr.AddPrivacyZone.Execute(cmd)
	require.NoError(t, err)

	require.False(t, privacyZoneUUID.IsZero())

	// test get
	_, err = tr.GetPrivacyZone.Execute(
		tracker.GetPrivacyZone{
			PrivacyZoneUUID: privacyZoneUUID,
			AsUser: &auth.ReadUser{
				UUID: user.UUID(),
			},
		},
	)
	require.NoError(t, err)

	// test list
	zones, err := tr.ListUserPrivacyZones.Execute(
		tracker.ListUserPrivacyZones{
			UserUUID: user.UUID(),
			AsUser: &auth.ReadUser{
				UUID: user.UUID(),
			},
		},
	)
	require.NoError(t, err)
	require.Len(t, zones, 1)

	// test delete
	err = tr.DeletePrivacyZone.Execute(
		tracker.DeletePrivacyZone{
			PrivacyZoneUUID: privacyZoneUUID,
			AsUser: &auth.ReadUser{
				UUID: user.UUID(),
			},
		},
	)
	require.NoError(t, err)

	// test get
	_, err = tr.GetPrivacyZone.Execute(
		tracker.GetPrivacyZone{
			PrivacyZoneUUID: privacyZoneUUID,
			AsUser: &auth.ReadUser{
				UUID: user.UUID(),
			},
		},
	)
	require.ErrorIs(t, err, tracker.ErrPrivacyZoneNotFound)

	// test list
	zones, err = tr.ListUserPrivacyZones.Execute(
		tracker.ListUserPrivacyZones{
			UserUUID: user.UUID(),
			AsUser: &auth.ReadUser{
				UUID: user.UUID(),
			},
		},
	)
	require.NoError(t, err)
	require.Empty(t, zones)
}

func TestDeletePrivacyZonePermissions(t *testing.T) {
	user := auth.ReadUser{
		UUID: auth.MustNewUserUUID("user-uuid"),
	}

	otherUser := auth.ReadUser{
		UUID: auth.MustNewUserUUID("other-user-uuid"),
	}

	testCases := []struct {
		Name      string
		User      *auth.ReadUser
		CanDelete bool
	}{
		{
			Name:      "unauthorized_user",
			User:      nil,
			CanDelete: false,
		},
		{
			Name:      "other_user",
			User:      &otherUser,
			CanDelete: false,
		},
		{
			Name:      "user",
			User:      &user,
			CanDelete: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			testTracker, cleanupTracker := NewTracker(t)
			defer cleanupTracker()

			tr := testTracker.Tracker

			position := domain.NewPosition(
				domain.MustNewLatitude(48.20952),
				domain.MustNewLongitude(16.35618),
			)

			circle := domain.MustNewCircle(
				domain.NewPosition(
					domain.MustNewLatitude(48.21019),
					domain.MustNewLongitude(16.36163),
				),
				500,
			)

			name := domain.MustNewPrivacyZoneName("Privacy zone")

			cmd := tracker.AddPrivacyZone{
				UserUUID: user.UUID,
				Position: position,
				Circle:   circle,
				Name:     name,
			}

			privacyZoneUUID, err := tr.AddPrivacyZone.Execute(cmd)
			require.NoError(t, err)

			err = tr.DeletePrivacyZone.Execute(
				tracker.DeletePrivacyZone{
					PrivacyZoneUUID: privacyZoneUUID,
					AsUser:          testCase.User,
				},
			)

			if testCase.CanDelete {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tracker.ErrDeletingPrivacyZoneForbidden)
			}
		})
	}
}

func TestPrivacyZonesPermissions(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	otherUser := auth.ReadUser{
		UUID: auth.MustNewUserUUID("other-user-uuid"),
	}

	user := mockUser()
	testTracker.UserRepository.Users[user.UUID()] = user
	readUser := user.AsReadUser()

	position := domain.NewPosition(
		domain.MustNewLatitude(48.20952),
		domain.MustNewLongitude(16.35618),
	)

	circle := domain.MustNewCircle(
		domain.NewPosition(
			domain.MustNewLatitude(48.21019),
			domain.MustNewLongitude(16.36163),
		),
		500,
	)

	name := domain.MustNewPrivacyZoneName("Privacy zone")

	cmd := tracker.AddPrivacyZone{
		UserUUID: user.UUID(),
		Position: position,
		Circle:   circle,
		Name:     name,
	}

	privacyZoneUUID, err := tr.AddPrivacyZone.Execute(cmd)
	require.NoError(t, err)

	t.Run("get unauthorised", func(t *testing.T) {
		_, err = tr.GetPrivacyZone.Execute(
			tracker.GetPrivacyZone{
				PrivacyZoneUUID: privacyZoneUUID,
				AsUser:          nil,
			},
		)

		require.ErrorIs(t, err, tracker.ErrGettingPrivacyZoneForbidden)
	})

	t.Run("get other", func(t *testing.T) {
		_, err = tr.GetPrivacyZone.Execute(
			tracker.GetPrivacyZone{
				PrivacyZoneUUID: privacyZoneUUID,
				AsUser:          &otherUser,
			},
		)

		require.ErrorIs(t, err, tracker.ErrGettingPrivacyZoneForbidden)
	})

	t.Run("get owner", func(t *testing.T) {
		_, err = tr.GetPrivacyZone.Execute(
			tracker.GetPrivacyZone{
				PrivacyZoneUUID: privacyZoneUUID,
				AsUser:          &readUser,
			},
		)

		require.NoError(t, err)
	})

	t.Run("list unauthorised", func(t *testing.T) {
		_, err := tr.ListUserPrivacyZones.Execute(
			tracker.ListUserPrivacyZones{
				UserUUID: user.UUID(),
				AsUser:   nil,
			},
		)

		require.ErrorIs(t, err, tracker.ErrGettingPrivacyZoneForbidden)
	})

	t.Run("list other", func(t *testing.T) {
		_, err := tr.ListUserPrivacyZones.Execute(
			tracker.ListUserPrivacyZones{
				UserUUID: user.UUID(),
				AsUser:   &otherUser,
			},
		)

		require.ErrorIs(t, err, tracker.ErrGettingPrivacyZoneForbidden)
	})

	t.Run("list owner", func(t *testing.T) {
		result, err := tr.ListUserPrivacyZones.Execute(
			tracker.ListUserPrivacyZones{
				UserUUID: user.UUID(),
				AsUser:   &readUser,
			},
		)

		require.NoError(t, err)
		require.NotEmpty(t, result)
	})
}

func TestStravaExport(t *testing.T) {
	testCases := []struct {
		File string
		Len  int
	}{
		{
			File: testStravaExportFile,
			Len:  2,
		},
		{
			File: testStravaExportFileWithFit,
			Len:  1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.File, func(t *testing.T) {
			testTracker, cleanupTracker := NewTracker(t)
			defer cleanupTracker()

			tr := testTracker.Tracker

			stravaExportFile, cleanupFile := fixture.TestDataFile(t, testCase.File)
			defer cleanupFile()

			fi, err := stravaExportFile.Stat()
			require.NoError(t, err)

			user := mockUser()
			testTracker.UserRepository.Users[user.UUID()] = user
			readUser := user.AsReadUser()

			cmd := tracker.ImportStrava{
				Archive:     stravaExportFile,
				ArchiveSize: fi.Size(),
				UserUUID:    user.UUID(),
			}

			err = tr.ImportStrava.Execute(cmd)
			require.NoError(t, err)

			// test list
			activities, err := tr.ListUserActivities.Execute(
				tracker.ListUserActivities{
					UserUUID: user.UUID(),
					AsUser:   &readUser,
				},
			)
			require.NoError(t, err)

			require.Len(t, activities.Activities, testCase.Len)
			for _, activity := range activities.Activities {
				require.Equal(t,
					domain.PrivateActivityVisibility,
					activity.Activity.Visibility(),
				)
			}
		})
	}
}

func TestStravaExportWithFit(t *testing.T) {
	testTracker, cleanupTracker := NewTracker(t)
	defer cleanupTracker()

	tr := testTracker.Tracker

	stravaExportFile, cleanupFile := fixture.TestDataFile(t, testStravaExportFileWithFit)
	defer cleanupFile()

	fi, err := stravaExportFile.Stat()
	require.NoError(t, err)

	user := mockUser()
	testTracker.UserRepository.Users[user.UUID()] = user
	readUser := user.AsReadUser()

	cmd := tracker.ImportStrava{
		Archive:     stravaExportFile,
		ArchiveSize: fi.Size(),
		UserUUID:    user.UUID(),
	}

	err = tr.ImportStrava.Execute(cmd)
	require.NoError(t, err)

	// test list
	activities, err := tr.ListUserActivities.Execute(
		tracker.ListUserActivities{
			UserUUID: user.UUID(),
			AsUser:   &readUser,
		},
	)
	require.NoError(t, err)

	require.NotEmpty(t, activities.Activities)
	for _, activity := range activities.Activities {
		require.Equal(t,
			domain.PrivateActivityVisibility,
			activity.Activity.Visibility(),
		)
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

func mockUser() auth.User {
	return auth.MustNewUser(
		auth.MustNewUserUUID("user-uuid"),
		auth.MustNewUsername("username"),
		auth.MustNewDisplayName("display-name"),
		auth.PasswordHash("made-up"),
		false,
	)
}
