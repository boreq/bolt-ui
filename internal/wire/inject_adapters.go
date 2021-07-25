package wire

import (
	"github.com/boreq/velo/adapters"
	"github.com/boreq/velo/application"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

//lint:ignore U1000 because
var adaptersSet = wire.NewSet(
	//trackerAdapters.NewRouteFileParser,
	//wire.Bind(new(tracker.RouteFileParser), new(*trackerAdapters.RouteFileParser)),

	//trackerAdapters.NewRouteFileParserFit,
	//trackerAdapters.NewRouteFileParserGpx,

	//trackerAdapters.NewStravaExportFileParser,
	//wire.Bind(new(tracker.StravaExportFileParser), new(*trackerAdapters.StravaExportFileParser)),

	adapters.NewTransactionProvider,
	wire.Bind(new(application.TransactionProvider), new(*adapters.TransactionProvider)),

	newAdaptersProvider,
	wire.Bind(new(adapters.AdaptersProvider), new(*adaptersProvider)),
)

//lint:ignore U1000 because
var testAdaptersSet = wire.NewSet(
	//trackerAdapters.NewRouteFileParser,
	//wire.Bind(new(tracker.RouteFileParser), new(*trackerAdapters.RouteFileParser)),

	//trackerAdapters.NewRouteFileParserFit,
	//trackerAdapters.NewRouteFileParserGpx,

	//trackerAdapters.NewStravaExportFileParser,
	//wire.Bind(new(tracker.StravaExportFileParser), new(*trackerAdapters.StravaExportFileParser)),

	adapters.NewTransactionProvider,
	wire.Bind(new(application.TransactionProvider), new(*adapters.TransactionProvider)),

	newAdaptersProvider,
	wire.Bind(new(adapters.AdaptersProvider), new(*adaptersProvider)),
)

//lint:ignore U1000 because
var transactableAdaptersSet = wire.NewSet(
	wire.Struct(new(application.TransactableAdapters), "*"),

	//trackerAdapters.NewActivityRepository,
	//wire.Bind(new(tracker.ActivityRepository), new(*trackerAdapters.ActivityRepository)),

	//trackerAdapters.NewRouteRepository,
	//wire.Bind(new(tracker.RouteRepository), new(*trackerAdapters.RouteRepository)),

	//trackerAdapters.NewUserToActivityRepository,
	//wire.Bind(new(tracker.UserToActivityRepository), new(*trackerAdapters.UserToActivityRepository)),

	//trackerAdapters.NewPrivacyZoneRepository,
	//wire.Bind(new(tracker.PrivacyZoneRepository), new(*trackerAdapters.PrivacyZoneRepository)),

	//trackerAdapters.NewUserToPrivacyZoneRepository,
	//wire.Bind(new(tracker.UserToPrivacyZoneRepository), new(*trackerAdapters.UserToPrivacyZoneRepository)),
)

//lint:ignore U1000 because
var testTransactableAdaptersSet = wire.NewSet(
	wire.Struct(new(application.TransactableAdapters), "*"),

	//trackerAdapters.NewActivityRepository,
	//wire.Bind(new(tracker.ActivityRepository), new(*trackerAdapters.ActivityRepository)),

	//trackerAdapters.NewRouteRepository,
	//wire.Bind(new(tracker.RouteRepository), new(*trackerAdapters.RouteRepository)),

	//trackerAdapters.NewUserToActivityRepository,
	//wire.Bind(new(tracker.UserToActivityRepository), new(*trackerAdapters.UserToActivityRepository)),

	//trackerAdapters.NewPrivacyZoneRepository,
	//wire.Bind(new(tracker.PrivacyZoneRepository), new(*trackerAdapters.PrivacyZoneRepository)),

	//trackerAdapters.NewUserToPrivacyZoneRepository,
	//wire.Bind(new(tracker.UserToPrivacyZoneRepository), new(*trackerAdapters.UserToPrivacyZoneRepository)),
)

type adaptersProvider struct {
}

func newAdaptersProvider() *adaptersProvider {
	return &adaptersProvider{}
}

func (p *adaptersProvider) Provide(tx *bolt.Tx) (*application.TransactableAdapters, error) {
	return BuildTransactableAdapters(tx)
}

type testAdaptersProvider struct {
	mocks Mocks
}

func newTestAdaptersProvider(mocks Mocks) *testAdaptersProvider {
	return &testAdaptersProvider{
		mocks: mocks,
	}
}

func (p *testAdaptersProvider) Provide(tx *bolt.Tx) (*application.TransactableAdapters, error) {
	return BuildTestTransactableAdapters(tx, p.mocks)
}
