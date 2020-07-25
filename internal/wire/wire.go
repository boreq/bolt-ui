//+build wireinject

package wire

import (
	"github.com/boreq/velo/application/auth"
	"github.com/boreq/velo/application/queries"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/internal/config"
	"github.com/boreq/velo/internal/service"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

func BuildTransactableAuthRepositories(tx *bolt.Tx) (*auth.TransactableRepositories, error) {
	wire.Build(
		appSet,
	)

	return nil, nil
}

func BuildTransactableQueryRepositories(tx *bolt.Tx) (*queries.TransactableRepositories, error) {
	wire.Build(
		appSet,
	)

	return nil, nil
}

func BuildTransactableTrackerRepositories(tx *bolt.Tx) (*tracker.TransactableRepositories, error) {
	wire.Build(
		trackerSet,
	)

	return nil, nil
}

func BuildTrackerForTest(db *bolt.DB) (*tracker.Tracker, error) {
	wire.Build(
		trackerSet,
	)

	return nil, nil
}

func BuildAuthForTest(db *bolt.DB) (*auth.Auth, error) {
	wire.Build(
		appSet,
	)

	return nil, nil
}

func BuildAuth(conf *config.Config) (*auth.Auth, error) {
	wire.Build(
		appSet,
		boltSet,
	)

	return nil, nil
}

func BuildService(conf *config.Config) (*service.Service, error) {
	wire.Build(
		service.NewService,
		httpSet,
		appSet,
		musicSet,
		boltSet,
	)

	return nil, nil
}
