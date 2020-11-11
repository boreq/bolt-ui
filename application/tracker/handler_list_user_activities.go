package tracker

import (
	"sort"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
)

const activitiesPerPage = 10

type ListUserActivitiesResult struct {
	Activities []ActivityWithRoute
	HasPrev    bool
	HasNext    bool
}

type ListUserActivities struct {
	UserUUID    auth.UserUUID
	StartAfter  *domain.ActivityUUID
	StartBefore *domain.ActivityUUID
}

type ListUserActivitiesHandler struct {
	transactionProvider TransactionProvider
}

func NewListUserActivitiesHandler(transactionProvider TransactionProvider) *ListUserActivitiesHandler {
	return &ListUserActivitiesHandler{
		transactionProvider: transactionProvider,
	}
}

func (h *ListUserActivitiesHandler) Execute(query ListUserActivities) (ListUserActivitiesResult, error) {
	var result ListUserActivitiesResult

	listFn, err := h.getIteratorFunc(query)
	if err != nil {
		return ListUserActivitiesResult{}, errors.Wrap(err, "could not get a list function")
	}

	if err := h.transactionProvider.Read(func(adapters *TransactableRepositories) error {
		iter, err := listFn(adapters.UserToActivity, query.UserUUID)
		if err != nil {
			return errors.Wrap(err, "could not get the iterator")
		}

		for i := 0; i < activitiesPerPage; i++ {
			activity, ok := iter.Next()
			if !ok {
				break
			}

			route, err := adapters.Route.Get(activity.RouteUUID())
			if err != nil {
				return errors.Wrap(err, "could not get a route")
			}

			result.Activities = append(result.Activities, ActivityWithRoute{
				Activity: activity,
				Route:    route,
			})
		}

		result.HasPrev, result.HasNext = h.getPrevNext(query, iter)

		if err := iter.Error(); err != nil {
			return errors.Wrap(err, "iterator error")
		}

		return nil
	}); err != nil {
		return ListUserActivitiesResult{}, errors.Wrap(err, "transaction failed")
	}

	sort.Slice(result.Activities, func(i, j int) bool {
		return result.Activities[i].Activity.UUID().String() > result.Activities[j].Activity.UUID().String()
	})

	return result, nil
}

func (h *ListUserActivitiesHandler) getPrevNext(query ListUserActivities, iter ActivityIterator) (bool, bool) {
	if query.StartAfter != nil {
		_, hasNext := iter.Next()
		return true, hasNext
	}

	if query.StartBefore != nil {
		_, hasPrev := iter.Next()
		return hasPrev, true
	}

	_, hasNext := iter.Next()
	return false, hasNext
}

type listFn func(r UserToActivityRepository, userUUID auth.UserUUID) (ActivityIterator, error)

func (h *ListUserActivitiesHandler) getIteratorFunc(query ListUserActivities) (listFn, error) {
	if query.StartAfter != nil && query.StartBefore != nil {
		return nil, errors.New("specified after and before at the same time")
	}

	if query.StartAfter != nil {
		return func(r UserToActivityRepository, userUUID auth.UserUUID) (ActivityIterator, error) {
			return r.ListAfter(userUUID, *query.StartAfter)
		}, nil
	}

	if query.StartBefore != nil {
		return func(r UserToActivityRepository, userUUID auth.UserUUID) (ActivityIterator, error) {
			return r.ListBefore(userUUID, *query.StartBefore)
		}, nil
	}

	return func(r UserToActivityRepository, userUUID auth.UserUUID) (ActivityIterator, error) {
		return r.List(userUUID)
	}, nil
}
