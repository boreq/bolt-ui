package auth

import (
	"github.com/boreq/velo/domain/auth"
	"github.com/pkg/errors"
)

type Logout struct {
	Token auth.AccessToken
}

type LogoutHandler struct {
	transactionProvider  TransactionProvider
	accessTokenGenerator AccessTokenGenerator
}

func NewLogoutHandler(
	transactionProvider TransactionProvider,
	accessTokenGenerator AccessTokenGenerator,
) *LogoutHandler {
	return &LogoutHandler{
		transactionProvider:  transactionProvider,
		accessTokenGenerator: accessTokenGenerator,
	}
}

func (h *LogoutHandler) Execute(cmd Logout) error {
	username, err := h.accessTokenGenerator.GetUsername(cmd.Token)
	if err != nil {
		return errors.Wrap(err, "could not extract the username")
	}

	if err := h.transactionProvider.Write(func(r *TransactableRepositories) error {
		u, err := r.Users.Get(username)
		if err != nil {
			return errors.Wrap(err, "could not get the user")
		}

		if err := u.Logout(cmd.Token); err != nil {
			return errors.Wrap(err, "could not log out the user")
		}

		return r.Users.Put(*u)
	}); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}
