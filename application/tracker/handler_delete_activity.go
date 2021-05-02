package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/domain/permissions"
)

type DeleteActivity struct {
	ActivityUUID domain.ActivityUUID
	AsUser       *auth.ReadUser
}

type DeleteActivityHandler struct {
	transactionProvider TransactionProvider
}

func NewDeleteActivityHandler(transactionProvider TransactionProvider) *DeleteActivityHandler {
	return &DeleteActivityHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *DeleteActivityHandler) Execute(cmd DeleteActivity) error {
	if err := h.transactionProvider.Write(func(repositories *TransactableRepositories) error {
		activity, err := repositories.Activity.Get(cmd.ActivityUUID)
		if err != nil {
			return errors.Wrap(err, "could not get an activity")
		}

		ok, err := permissions.CanDeleteActivity(activity, cmd.AsUser)
		if err != nil {
			return errors.Wrap(err, "permission check failed")
		}

		if !ok {
			return ErrDeletingActivityForbidden
		}

		if err := repositories.UserToActivity.Unassign(activity.UserUUID(), activity.UUID()); err != nil {
			return errors.Wrap(err, "could not unassign the activity")
		}

		if err := repositories.Activity.Delete(activity.UUID()); err != nil {
			return errors.Wrap(err, "could not delete the activity")
		}

		if err := repositories.Route.Delete(activity.RouteUUID()); err != nil {
			return errors.Wrap(err, "could not delete the route")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}
