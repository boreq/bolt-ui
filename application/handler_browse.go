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

func (h *BrowseHandler) Execute(query Browse) (entries []Entry, err error) {
	if query.Before != nil && query.After != nil {
		return entries, errors.New("passed both before and after")
	}

	if err := h.transactionProvider.Read(func(adapters *TransactableAdapters) error {
		entries, err = adapters.Database.Browse(query.Path, query.Before, query.After)
		if err != nil {
			return errors.Wrap(err, "could not browse the database")
		}

		return nil
	}); err != nil {
		return entries, errors.Wrap(err, "transaction failed")
	}

	return entries, nil
}
