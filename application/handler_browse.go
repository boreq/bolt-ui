package application

import (
	"github.com/boreq/errors"
)

type Browse struct {
	Path   []Key
	Before *Key
	After  *Key
	From   *Key
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
	if !h.queryValid(query) {
		return tree, errors.New("passed two or more of before/after/from at the same time")
	}

	tree.Path = query.Path

	if err := h.transactionProvider.Read(func(adapters *TransactableAdapters) error {
		tree.Entries, err = adapters.Database.Browse(query.Path, query.Before, query.After, query.From)
		if err != nil {
			return errors.Wrap(err, "could not browse the database")
		}

		return nil
	}); err != nil {
		return tree, errors.Wrap(err, "transaction failed")
	}

	return tree, nil
}

func (h *BrowseHandler) queryValid(query Browse) bool {
	var counter int
	if query.Before != nil {
		counter++
	}
	if query.After != nil {
		counter++
	}
	if query.From != nil {
		counter++
	}
	return counter <= 1
}
