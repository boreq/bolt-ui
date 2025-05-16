package wire

import (
	bolt2 "github.com/boreq/bolt-ui/adapters/bolt"
	"github.com/boreq/bolt-ui/application"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

//lint:ignore U1000 because
var adaptersSet = wire.NewSet(
	bolt2.NewTransactionProvider,
	wire.Bind(new(application.TransactionProvider), new(*bolt2.TransactionProvider)),

	newAdaptersProvider,
	wire.Bind(new(bolt2.AdaptersProvider), new(*adaptersProvider)),
)

//lint:ignore U1000 because
var testAdaptersSet = wire.NewSet(
	bolt2.NewTransactionProvider,
	wire.Bind(new(application.TransactionProvider), new(*bolt2.TransactionProvider)),

	newTestAdaptersProvider,
	wire.Bind(new(bolt2.AdaptersProvider), new(*testAdaptersProvider)),
)

//lint:ignore U1000 because
var transactableAdaptersSet = wire.NewSet(
	wire.Struct(new(application.TransactableAdapters), "*"),

	bolt2.NewDatabase,
	wire.Bind(new(application.Database), new(*bolt2.Database)),
)

//lint:ignore U1000 because
var testTransactableAdaptersSet = wire.NewSet(
	wire.Struct(new(application.TransactableAdapters), "*"),

	bolt2.NewDatabase,
	wire.Bind(new(application.Database), new(*bolt2.Database)),
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
