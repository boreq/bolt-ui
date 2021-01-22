package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/domain/permissions"
)

type EditActivity struct {
	ActivityUUID domain.ActivityUUID
	AsUser       *auth.ReadUser

	Title      domain.ActivityTitle
	Visibility domain.ActivityVisibility
}

type EditActivityHandler struct {
	transactionProvider TransactionProvider
}

func NewEditActivityHandler(transactionProvider TransactionProvider) *EditActivityHandler {
	return &EditActivityHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *EditActivityHandler) Execute(cmd EditActivity) error {
	if err := h.transactionProvider.Write(func(repositories *TransactableRepositories) error {
		activity, err := repositories.Activity.Get(cmd.ActivityUUID)
		if err != nil {
			return errors.Wrap(err, "could not get an activity")
		}

		ok, err := permissions.CanEditActivity(activity, cmd.AsUser)
		if err != nil {
			return errors.Wrap(err, "permission check failed")
		}

		if !ok {
			return ErrEditingActivityForbidden
		}

		if err := activity.ChangeTitle(cmd.Title); err != nil {
			return errors.Wrap(err, "could not change the title")
		}

		if err := activity.ChangeVisibility(cmd.Visibility); err != nil {
			return errors.Wrap(err, "could not change the visibility")
		}

		if !activity.HasChanges() {
			return nil
		}

		if err := repositories.Activity.Save(activity); err != nil {
			return errors.Wrap(err, "could not save an activity")
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}
