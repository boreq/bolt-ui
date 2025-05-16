package application

import (
	"github.com/boreq/errors"
)

type Browse struct {
	path   []Key
	before *Key
	after  *Key
	from   *Key
}

func NewBrowse(path []Key, before *Key, after *Key, from *Key) (Browse, error) {
	var counter int
	if before != nil {
		counter++
	}
	if after != nil {
		counter++
	}
	if from != nil {
		counter++
	}
	if counter > 1 {
		return Browse{}, errors.New("passed more than one before/after/from at the same time")
	}

	return Browse{
		path:   path,
		before: before,
		after:  after,
		from:   from,
	}, nil
}

func (b Browse) Path() []Key {
	return b.path
}

func (b Browse) Before() *Key {
	return b.before
}

func (b Browse) After() *Key {
	return b.after
}

func (b Browse) From() *Key {
	return b.from
}

type BrowseHandler struct {
	transactionProvider TransactionProvider
}

func NewBrowseHandler(transactionProvider TransactionProvider) *BrowseHandler {
	return &BrowseHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *BrowseHandler) Execute(query Browse) (tree Tree, err error) {
	tree.Path = query.Path()

	if err := h.transactionProvider.Read(func(adapters *TransactableAdapters) error {
		tree.Entries, err = adapters.Database.Browse(query.Path(), query.Before(), query.After(), query.From())
		if err != nil {
			return errors.Wrap(err, "could not browse the database")
		}

		return nil
	}); err != nil {
		return tree, errors.Wrap(err, "transaction failed")
	}

	return tree, nil
}
