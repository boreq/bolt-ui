package auth

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/domain/permissions"
)

var ErrChangingPasswordForbidden = errors.New("this user can not change this password")

type ChangePassword struct {
	Username    string
	OldPassword string
	NewPassword auth.ValidatedPassword
	AsUser      *auth.ReadUser
}

type ChangePasswordHandler struct {
	transactionProvider TransactionProvider
	passwordHasher      PasswordHasher
}

func NewChangePasswordHandler(
	transactionProvider TransactionProvider,
	passwordHasher PasswordHasher,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		transactionProvider: transactionProvider,
		passwordHasher:      passwordHasher,
	}
}

func (h *ChangePasswordHandler) Execute(cmd ChangePassword) error {
	if cmd.NewPassword.IsZero() {
		return errors.New("zero value of new password")
	}

	passwordHash, err := h.passwordHasher.Hash(cmd.NewPassword.String())
	if err != nil {
		return errors.Wrap(err, "hashing the password failed")
	}

	return h.transactionProvider.Write(func(r *TransactableRepositories) error {
		u, err := r.Users.Get(cmd.Username)
		if err != nil {
			return errors.Wrap(err, "could not get the user")
		}

		ok, err := permissions.CanChangePassword(toReadUser(*u), cmd.AsUser)
		if err != nil {
			return errors.Wrap(err, "permission check failed")
		}

		if !ok {
			return ErrChangingPasswordForbidden
		}

		if err := h.passwordHasher.Compare(u.Password, cmd.OldPassword); err != nil {
			return errors.Wrap(ErrUnauthorized, "invalid password")
		}

		u.Password = passwordHash

		return r.Users.Put(*u)
	})
}
