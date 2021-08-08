//+build wireinject

package wire

import (
	"github.com/boreq/bolt-ui/application"
	"github.com/boreq/bolt-ui/internal/config"
	"github.com/boreq/bolt-ui/internal/service"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

func BuildTransactableAdapters(_ *bolt.Tx) (*application.TransactableAdapters, error) {
	wire.Build(
		transactableAdaptersSet,
	)

	return nil, nil
}

func BuildTestTransactableAdapters(_ *bolt.Tx, _ Mocks) (*application.TransactableAdapters, error) {
	wire.Build(
		testTransactableAdaptersSet,
	)

	return nil, nil
}

func BuildApplicationForTest(db *bolt.DB) (TestApplication, error) {
	wire.Build(
		appSet,
		testAdaptersSet,

		wire.Struct(new(TestApplication), "*"),
		wire.Struct(new(Mocks), "*"),
	)

	return TestApplication{}, nil
}

type TestApplication struct {
	Application *application.Application
	Mocks
	DB *bolt.DB
}

type Mocks struct {
}

func BuildService(conf *config.Config) (*service.Service, error) {
	wire.Build(
		service.NewService,
		httpSet,
		appSet,
		boltSet,
		adaptersSet,
	)

	return nil, nil
}
