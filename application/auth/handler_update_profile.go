package auth

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/domain/permissions"
)

var ErrUpdatingProfileForbidden = errors.New("this user can not update this profile")

type UpdateProfile struct {
	Username    string
	DisplayName auth.ValidatedDisplayName
	AsUser      *auth.ReadUser
}

type UpdateProfileHandler struct {
	transactionProvider TransactionProvider
}

func NewUpdateProfileHandler(
	transactionProvider TransactionProvider,
) *UpdateProfileHandler {
	return &UpdateProfileHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *UpdateProfileHandler) Execute(cmd UpdateProfile) error {
	if cmd.DisplayName.IsZero() {
		return errors.New("zero value of display name")
	}

	return h.transactionProvider.Write(func(r *TransactableRepositories) error {
		u, err := r.Users.Get(cmd.Username)
		if err != nil {
			return errors.Wrap(err, "could not get the user")
		}

		ok, err := permissions.CanUpdateProfile(toReadUser(*u), cmd.AsUser)
		if err != nil {
			return errors.Wrap(err, "permission check failed")
		}

		if !ok {
			return ErrUpdatingProfileForbidden
		}

		u.DisplayName = cmd.DisplayName.String()

		return r.Users.Put(*u)
	})
}
