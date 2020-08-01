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

	trackerAdapters.NewUUIDGenerator,
	wire.Bind(new(tracker.UUIDGenerator), new(*trackerAdapters.UUIDGenerator)),
)

//lint:ignore U1000 because
var trackerTransactableRepositoriesSet = wire.NewSet(
	trackerAdapters.NewTrackerTransactionProvider,
	wire.Bind(new(tracker.TransactionProvider), new(*trackerAdapters.TrackerTransactionProvider)),

	wire.Struct(new(tracker.TransactableRepositories), "*"),

	newTrackerRepositoriesProvider,
	wire.Bind(new(trackerAdapters.TrackerRepositoriesProvider), new(*trackerRepositoriesProvider)),

	trackerAdapters.NewActivityRepository,
	wire.Bind(new(tracker.ActivityRepository), new(*trackerAdapters.ActivityRepository)),

	trackerAdapters.NewRouteRepository,
	wire.Bind(new(tracker.RouteRepository), new(*trackerAdapters.RouteRepository)),
)

type trackerRepositoriesProvider struct {
}

func newTrackerRepositoriesProvider() *trackerRepositoriesProvider {
	return &trackerRepositoriesProvider{}
}

func (p *trackerRepositoriesProvider) Provide(tx *bolt.Tx) (*tracker.TransactableRepositories, error) {
	return BuildTransactableTrackerRepositories(tx)
}
