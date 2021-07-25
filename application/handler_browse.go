package application

import (
	"github.com/boreq/errors"
)

type Browse struct {
}

type BrowseHandler struct {
	transactionProvider TransactionProvider
}

func NewBrowseHandler(transactionProvider TransactionProvider) *BrowseHandler {
	return &BrowseHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *BrowseHandler) Execute(query Browse) error {
	if err := h.transactionProvider.Read(func(adapters *TransactableAdapters) error {
		return errors.New("not implemented")
	}); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}
