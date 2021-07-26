package application

import (
	"github.com/boreq/errors"
)

type Browse struct {
	Path   []Key
	Before *Key
	After  *Key
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
	if query.Before != nil && query.After != nil {
		return tree, errors.New("passed both before and after")
	}

	tree.Path = query.Path

	if err := h.transactionProvider.Read(func(adapters *TransactableAdapters) error {
		tree.Entries, err = adapters.Database.Browse(query.Path, query.Before, query.After)
		if err != nil {
			return errors.Wrap(err, "could not browse the database")
		}

		return nil
	}); err != nil {
		return tree, errors.Wrap(err, "transaction failed")
	}

	return tree, nil
}
