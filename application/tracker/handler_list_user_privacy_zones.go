package tracker

import (
	"sort"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/domain/permissions"
)

type ListUserPrivacyZones struct {
	UserUUID auth.UserUUID
	AsUser   *auth.ReadUser
}

type ListUserPrivacyZonesHandler struct {
	transactionProvider TransactionProvider
}

func NewListUserPrivacyZonesHandler(transactionProvider TransactionProvider) *ListUserPrivacyZonesHandler {
	return &ListUserPrivacyZonesHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *ListUserPrivacyZonesHandler) Execute(query ListUserPrivacyZones) ([]*domain.PrivacyZone, error) {
	if !permissions.CanListPrivacyZones(query.UserUUID, query.AsUser) {
		return nil, ErrGettingPrivacyZoneForbidden
	}

	var result []*domain.PrivacyZone

	if err := h.transactionProvider.Read(func(adapters *TransactableRepositories) error {
		privacyZones, err := adapters.UserToPrivacyZone.List(query.UserUUID)
		if err != nil {
			return errors.Wrap(err, "could not list privacy zones")
		}

		for _, privacyZone := range privacyZones {
			ok, err := permissions.CanViewPrivacyZone(privacyZone, query.AsUser)
			if err != nil {
				return errors.Wrap(err, "permission check failed")
			}

			if ok {
				result = append(result, privacyZone)
			}
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].UUID().String() > result[j].UUID().String()
	})

	return result, nil
}
