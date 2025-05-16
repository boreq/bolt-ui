package wire

import (
	boltadapters "github.com/boreq/bolt-ui/adapters/bolt"
	"github.com/boreq/bolt-ui/application"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

//lint:ignore U1000 because
var adaptersSet = wire.NewSet(
	boltadapters.NewTransactionProvider,
	wire.Bind(new(application.TransactionProvider), new(*boltadapters.TransactionProvider)),

	newAdaptersProvider,
	wire.Bind(new(boltadapters.AdaptersProvider), new(*adaptersProvider)),
)

//lint:ignore U1000 because
var testAdaptersSet = wire.NewSet(
	boltadapters.NewTransactionProvider,
	wire.Bind(new(application.TransactionProvider), new(*boltadapters.TransactionProvider)),

	newTestAdaptersProvider,
	wire.Bind(new(boltadapters.AdaptersProvider), new(*testAdaptersProvider)),
)

//lint:ignore U1000 because
var transactableAdaptersSet = wire.NewSet(
	wire.Struct(new(application.TransactableAdapters), "*"),

	boltadapters.NewDatabase,
	wire.Bind(new(application.Database), new(*boltadapters.Database)),
)

//lint:ignore U1000 because
var testTransactableAdaptersSet = wire.NewSet(
	wire.Struct(new(application.TransactableAdapters), "*"),

	boltadapters.NewDatabase,
	wire.Bind(new(application.Database), new(*boltadapters.Database)),
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
