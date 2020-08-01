package auth

import "github.com/boreq/errors"

type GetUser struct {
	Username string
}

type GetUserHandler struct {
	transactionProvider TransactionProvider
}

func NewGetUserHandler(transactionProvider TransactionProvider) *GetUserHandler {
	return &GetUserHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *GetUserHandler) Execute(query GetUser) (ReadUser, error) {
	var user User
	if err := h.transactionProvider.Read(func(r *TransactableRepositories) error {
		u, err := r.Users.Get(query.Username)
		if err != nil {
			return errors.Wrap(err, "could not get the user")
		}
		user = *u
		return nil
	}); err != nil {
		return ReadUser{}, errors.Wrap(err, "transaction failed")
	}
	return toReadUser(user), nil
}
