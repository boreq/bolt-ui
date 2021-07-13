package wire

import (
	trackerAdapters "github.com/boreq/velo/adapters/tracker"
	"github.com/boreq/velo/application/tracker"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

//lint:ignore U1000 because
var trackerRepositoriesSet = wire.NewSet(
	trackerAdapters.NewRouteFileParser,
	wire.Bind(new(tracker.RouteFileParser), new(*trackerAdapters.RouteFileParser)),

	trackerAdapters.NewRouteFileParserFit,
	trackerAdapters.NewRouteFileParserGpx,

	trackerAdapters.NewStravaExportFileParser,
	wire.Bind(new(tracker.StravaExportFileParser), new(*trackerAdapters.StravaExportFileParser)),

	trackerAdapters.NewTrackerTransactionProvider,
	wire.Bind(new(tracker.TransactionProvider), new(*trackerAdapters.TrackerTransactionProvider)),

	newTrackerRepositoriesProvider,
	wire.Bind(new(trackerAdapters.TrackerRepositoriesProvider), new(*trackerRepositoriesProvider)),
)

//lint:ignore U1000 because
var trackerTestRepositoriesSet = wire.NewSet(
	trackerAdapters.NewRouteFileParser,
	wire.Bind(new(tracker.RouteFileParser), new(*trackerAdapters.RouteFileParser)),

	trackerAdapters.NewRouteFileParserFit,
	trackerAdapters.NewRouteFileParserGpx,

	trackerAdapters.NewStravaExportFileParser,
	wire.Bind(new(tracker.StravaExportFileParser), new(*trackerAdapters.StravaExportFileParser)),

	trackerAdapters.NewTrackerTransactionProvider,
	wire.Bind(new(tracker.TransactionProvider), new(*trackerAdapters.TrackerTransactionProvider)),

	newTrackerTestRepositoriesProvider,
	wire.Bind(new(trackerAdapters.TrackerRepositoriesProvider), new(*trackerTestRepositoriesProvider)),
)

//lint:ignore U1000 because
var trackerTransactableRepositoriesSet = wire.NewSet(
	wire.Struct(new(tracker.TransactableRepositories), "*"),

	trackerAdapters.NewActivityRepository,
	wire.Bind(new(tracker.ActivityRepository), new(*trackerAdapters.ActivityRepository)),

	trackerAdapters.NewRouteRepository,
	wire.Bind(new(tracker.RouteRepository), new(*trackerAdapters.RouteRepository)),

	trackerAdapters.NewUserToActivityRepository,
	wire.Bind(new(tracker.UserToActivityRepository), new(*trackerAdapters.UserToActivityRepository)),

	trackerAdapters.NewPrivacyZoneRepository,
	wire.Bind(new(tracker.PrivacyZoneRepository), new(*trackerAdapters.PrivacyZoneRepository)),

	trackerAdapters.NewUserToPrivacyZoneRepository,
	wire.Bind(new(tracker.UserToPrivacyZoneRepository), new(*trackerAdapters.UserToPrivacyZoneRepository)),
)

//lint:ignore U1000 because
var trackerTestTransactableRepositoriesSet = wire.NewSet(
	wire.Struct(new(tracker.TransactableRepositories), "*"),

	trackerAdapters.NewActivityRepository,
	wire.Bind(new(tracker.ActivityRepository), new(*trackerAdapters.ActivityRepository)),

	trackerAdapters.NewRouteRepository,
	wire.Bind(new(tracker.RouteRepository), new(*trackerAdapters.RouteRepository)),

	trackerAdapters.NewUserToActivityRepository,
	wire.Bind(new(tracker.UserToActivityRepository), new(*trackerAdapters.UserToActivityRepository)),

	trackerAdapters.NewPrivacyZoneRepository,
	wire.Bind(new(tracker.PrivacyZoneRepository), new(*trackerAdapters.PrivacyZoneRepository)),

	trackerAdapters.NewUserToPrivacyZoneRepository,
	wire.Bind(new(tracker.UserToPrivacyZoneRepository), new(*trackerAdapters.UserToPrivacyZoneRepository)),
)

type trackerRepositoriesProvider struct {
}

func newTrackerRepositoriesProvider() *trackerRepositoriesProvider {
	return &trackerRepositoriesProvider{}
}

func (p *trackerRepositoriesProvider) Provide(tx *bolt.Tx) (*tracker.TransactableRepositories, error) {
	return BuildTransactableTrackerRepositories(tx)
}

type trackerTestRepositoriesProvider struct {
	mocks TrackerMocks
}

func newTrackerTestRepositoriesProvider(mocks TrackerMocks) *trackerTestRepositoriesProvider {
	return &trackerTestRepositoriesProvider{
		mocks: mocks,
	}
}

func (p *trackerTestRepositoriesProvider) Provide(tx *bolt.Tx) (*tracker.TransactableRepositories, error) {
	return BuildTestTransactableTrackerRepositories(tx, p.mocks)
}
