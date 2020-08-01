package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
)

type GetActivityResult struct {
	Activity *domain.Activity
	Route    *domain.Route
}

type GetActivity struct {
	ActivityUUID domain.ActivityUUID
}

type GetActivityHandler struct {
	transactionProvider TransactionProvider
}

func NewGetActivityHandler(transactionProvider TransactionProvider) *GetActivityHandler {
	return &GetActivityHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *GetActivityHandler) Execute(query GetActivity) (GetActivityResult, error) {
	var result GetActivityResult

	if err := h.transactionProvider.Read(func(repositories *TransactableRepositories) error {
		activity, err := repositories.Activity.Get(query.ActivityUUID)
		if err != nil {
			return errors.Wrap(err, "could not get an activity")
		}

		route, err := repositories.Route.Get(activity.RouteUUID())
		if err != nil {
			return errors.Wrap(err, "could not get a route")
		}

		result = GetActivityResult{
			Activity: activity,
			Route:    route,
		}

		return nil
	}); err != nil {
		return result, errors.Wrap(err, "transaction failed")
	}

	return result, nil
}
