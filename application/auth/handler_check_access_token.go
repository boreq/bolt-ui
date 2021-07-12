package auth

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain/auth"
)

type CheckAccessToken struct {
	Token auth.AccessToken
}

type CheckAccessTokenHandler struct {
	transactionProvider  TransactionProvider
	accessTokenGenerator AccessTokenGenerator
}

func NewCheckAccessTokenHandler(
	transactionProvider TransactionProvider,
	accessTokenGenerator AccessTokenGenerator,
) *CheckAccessTokenHandler {
	return &CheckAccessTokenHandler{
		transactionProvider:  transactionProvider,
		accessTokenGenerator: accessTokenGenerator,
	}
}

func (h *CheckAccessTokenHandler) Execute(cmd CheckAccessToken) (*auth.ReadUser, error) {
	username, err := h.accessTokenGenerator.GetUsername(cmd.Token)
	if err != nil {
		return nil, errors.Wrap(ErrUnauthorized, "could not get the username")
	}

	var foundUser *auth.User
	if err := h.transactionProvider.Write(func(r *TransactableRepositories) error {
		u, err := r.Users.Get(username)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				return errors.Wrap(ErrUnauthorized, "user not found")
			}
			return errors.Wrap(err, "could not get the user")
		}

		ok, err := u.CheckAccessToken(cmd.Token)
		if err != nil {
			return errors.Wrap(err, "invalid access token")
		}

		if !ok {
			return errors.Wrap(ErrUnauthorized, "invalid token")
		}

		foundUser = u

		return r.Users.Put(*u)
	}); err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	rv := foundUser.AsReadUser()
	return &rv, nil
}
