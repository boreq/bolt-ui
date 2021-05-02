package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/domain/permissions"
)

type DeletePrivacyZone struct {
	PrivacyZoneUUID domain.PrivacyZoneUUID
	AsUser          *auth.ReadUser
}

type DeletePrivacyZoneHandler struct {
	transactionProvider TransactionProvider
}

func NewDeletePrivacyZoneHandler(transactionProvider TransactionProvider) *DeletePrivacyZoneHandler {
	return &DeletePrivacyZoneHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *DeletePrivacyZoneHandler) Execute(cmd DeletePrivacyZone) error {
	if err := h.transactionProvider.Write(func(repositories *TransactableRepositories) error {
		privacyZone, err := repositories.PrivacyZone.Get(cmd.PrivacyZoneUUID)
		if err != nil {
			return errors.Wrap(err, "could not get the privacy zone")
		}

		ok, err := permissions.CanDeletePrivacyZone(privacyZone, cmd.AsUser)
		if err != nil {
			return errors.Wrap(err, "permission check failed")
		}

		if !ok {
			return ErrDeletingPrivacyZoneForbidden
		}

		if err := repositories.UserToPrivacyZone.Unassign(privacyZone.UserUUID(), privacyZone.UUID()); err != nil {
			return errors.Wrap(err, "could not unassign the privacy zone")
		}

		if err := repositories.PrivacyZone.Delete(privacyZone.UUID()); err != nil {
			return errors.Wrap(err, "could not delete the privacy zone")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}
