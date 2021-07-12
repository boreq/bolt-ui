package auth

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain/auth"
)

type SetPassword struct {
	Username string
	Password auth.Password
}

type SetPasswordHandler struct {
	passwordHasher      PasswordHasher
	transactionProvider TransactionProvider
}

func NewSetPasswordHandler(
	passwordHasher PasswordHasher,
	transactionProvider TransactionProvider,
) *SetPasswordHandler {
	return &SetPasswordHandler{
		passwordHasher:      passwordHasher,
		transactionProvider: transactionProvider,
	}
}

func (h *SetPasswordHandler) Execute(cmd SetPassword) error {
	if cmd.Password.IsZero() {
		return errors.New("zero value of password")
	}

	passwordHash, err := h.passwordHasher.Hash(cmd.Password.String())
	if err != nil {
		return errors.Wrap(err, "hashing the password failed")
	}

	return h.transactionProvider.Write(func(r *TransactableRepositories) error {
		u, err := r.Users.Get(cmd.Username)
		if err != nil {
			return errors.Wrap(err, "could not get the user")
		}

		if err := u.ChangePassword(passwordHash); err != nil {
			return errors.Wrap(err, "could not change the password")
		}

		return r.Users.Put(*u)
	})
}
