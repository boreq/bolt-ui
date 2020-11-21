//+build wireinject

package wire

import (
	"github.com/boreq/velo/application/auth"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/internal/config"
	"github.com/boreq/velo/internal/service"
	"github.com/boreq/velo/internal/tests/mocks"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

func BuildTransactableAuthRepositories(tx *bolt.Tx) (*auth.TransactableRepositories, error) {
	wire.Build(
		appSet,
	)

	return nil, nil
}

func BuildTransactableTrackerRepositories(tx *bolt.Tx) (*tracker.TransactableRepositories, error) {
	wire.Build(
		trackerTransactableRepositoriesSet,
		authTransactableRepositoriesSet,
	)

	return nil, nil
}

func BuildTestTransactableTrackerRepositories(_ *bolt.Tx, _ TrackerMocks) (*tracker.TransactableRepositories, error) {
	wire.Build(
		trackerTestTransactableRepositoriesSet,

		wire.FieldsOf(new(TrackerMocks), "UserRepository"),
		wire.Bind(new(tracker.UserRepository), new(*mocks.TrackerUserRepositoryMock)),
	)

	return nil, nil
}

func BuildTrackerForTest(db *bolt.DB) (TestTracker, error) {
	wire.Build(
		trackerSet,
		trackerTestRepositoriesSet,

		mocks.NewTrackerUserRepositoryMock,

		adaptersSet,

		wire.Struct(new(TestTracker), "*"),
		wire.Struct(new(TrackerMocks), "*"),
	)

	return TestTracker{}, nil
}

type TestTracker struct {
	Tracker *tracker.Tracker
	TrackerMocks
}

type TrackerMocks struct {
	UserRepository *mocks.TrackerUserRepositoryMock
}

func BuildAuthForTest(db *bolt.DB) (*auth.Auth, error) {
	wire.Build(
		appSet,
		adaptersSet,
	)

	return nil, nil
}

func BuildAuth(conf *config.Config) (*auth.Auth, error) {
	wire.Build(
		appSet,
		boltSet,
		adaptersSet,
	)

	return nil, nil
}

func BuildService(conf *config.Config) (*service.Service, error) {
	wire.Build(
		service.NewService,
		httpSet,
		appSet,
		trackerRepositoriesSet,
		boltSet,
		adaptersSet,
	)

	return nil, nil
}
