package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/domain/permissions"
)

type GetPrivacyZone struct {
	PrivacyZoneUUID domain.PrivacyZoneUUID
	AsUser          *auth.ReadUser
}

type GetPrivacyZoneHandler struct {
	transactionProvider TransactionProvider
}

func NewGetPrivacyZoneHandler(transactionProvider TransactionProvider) *GetPrivacyZoneHandler {
	return &GetPrivacyZoneHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *GetPrivacyZoneHandler) Execute(query GetPrivacyZone) (*domain.PrivacyZone, error) {
	var result *domain.PrivacyZone

	if err := h.transactionProvider.Read(func(repositories *TransactableRepositories) error {
		privacyZone, err := repositories.PrivacyZone.Get(query.PrivacyZoneUUID)
		if err != nil {
			return errors.Wrap(err, "could not get an activity")
		}

		ok, err := permissions.CanViewPrivacyZone(privacyZone, query.AsUser)
		if err != nil {
			return errors.Wrap(err, "permission check failed")
		}

		if !ok {
			return ErrGettingPrivacyZoneForbidden
		}

		result = privacyZone

		return nil
	}); err != nil {
		return result, errors.Wrap(err, "transaction failed")
	}

	return result, nil
}
