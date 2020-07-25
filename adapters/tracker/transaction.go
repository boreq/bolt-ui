package tracker

import (
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/errors"
	bolt "go.etcd.io/bbolt"
)

type TrackerRepositoriesProvider interface {
	Provide(tx *bolt.Tx) (*tracker.TransactableRepositories, error)
}

type TrackerTransactionProvider struct {
	db                   *bolt.DB
	repositoriesProvider TrackerRepositoriesProvider
}

func NewTrackerTransactionProvider(
	db *bolt.DB,
	repositoriesProvider TrackerRepositoriesProvider,
) *TrackerTransactionProvider {
	return &TrackerTransactionProvider{
		db:                   db,
		repositoriesProvider: repositoriesProvider,
	}
}

func (p *TrackerTransactionProvider) Read(handler tracker.TransactionHandler) error {
	return p.db.View(func(tx *bolt.Tx) error {
		repositories, err := p.repositoriesProvider.Provide(tx)
		if err != nil {
			return errors.Wrap(err, "could not provide the repositories")
		}
		return handler(repositories)
	})
}

func (p *TrackerTransactionProvider) Write(handler tracker.TransactionHandler) error {
	return p.db.Update(func(tx *bolt.Tx) error {
		repositories, err := p.repositoriesProvider.Provide(tx)
		if err != nil {
			return errors.Wrap(err, "could not provide the repositories")
		}
		return handler(repositories)
	})
}
