package tracker

import "github.com/boreq/errors"

type AddActivity struct {
}

type AddActivityHandler struct {
	transactionProvider TransactionProvider
}

func NewAddActivityHandler(transactionProvider TransactionProvider) *AddActivityHandler {
	return &AddActivityHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *AddActivityHandler) Execute(cmd AddActivity) error {
	return errors.New("not implemented")
}
