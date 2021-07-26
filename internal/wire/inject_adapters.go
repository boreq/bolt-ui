package wire

import (
	"github.com/boreq/velo/adapters"
	"github.com/boreq/velo/application"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

//lint:ignore U1000 because
var adaptersSet = wire.NewSet(
	adapters.NewTransactionProvider,
	wire.Bind(new(application.TransactionProvider), new(*adapters.TransactionProvider)),

	newAdaptersProvider,
	wire.Bind(new(adapters.AdaptersProvider), new(*adaptersProvider)),
)

//lint:ignore U1000 because
var testAdaptersSet = wire.NewSet(
	adapters.NewTransactionProvider,
	wire.Bind(new(application.TransactionProvider), new(*adapters.TransactionProvider)),

	newTestAdaptersProvider,
	wire.Bind(new(adapters.AdaptersProvider), new(*testAdaptersProvider)),
)

//lint:ignore U1000 because
var transactableAdaptersSet = wire.NewSet(
	wire.Struct(new(application.TransactableAdapters), "*"),

	adapters.NewDatabase,
	wire.Bind(new(application.Database), new(*adapters.Database)),
)

//lint:ignore U1000 because
var testTransactableAdaptersSet = wire.NewSet(
	wire.Struct(new(application.TransactableAdapters), "*"),

	adapters.NewDatabase,
	wire.Bind(new(application.Database), new(*adapters.Database)),
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
