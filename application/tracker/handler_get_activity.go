package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/domain/permissions"
)

type GetActivity struct {
	ActivityUUID domain.ActivityUUID
	AsUser       *auth.ReadUser
}

type GetActivityHandler struct {
	transactionProvider TransactionProvider
}

func NewGetActivityHandler(transactionProvider TransactionProvider) *GetActivityHandler {
	return &GetActivityHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *GetActivityHandler) Execute(query GetActivity) (Activity, error) {
	var result Activity

	if err := h.transactionProvider.Read(func(repositories *TransactableRepositories) error {
		activity, err := repositories.Activity.Get(query.ActivityUUID)
		if err != nil {
			return errors.Wrap(err, "could not get an activity")
		}

		ok, err := permissions.CanViewActivity(activity, query.AsUser)
		if err != nil {
			return errors.Wrap(err, "permission check failed")
		}

		if !ok {
			return ErrGettingActivityForbidden
		}

		route, err := repositories.Route.Get(activity.RouteUUID())
		if err != nil {
			return errors.Wrap(err, "could not get a route")
		}

		user, err := repositories.User.GetByUUID(activity.UserUUID())
		if err != nil {
			return errors.Wrap(err, "could not get a user")
		}

		result = Activity{
			Activity: activity,
			Route:    route,
			User:     toUser(user),
		}

		return nil
	}); err != nil {
		return result, errors.Wrap(err, "transaction failed")
	}

	return result, nil
}
