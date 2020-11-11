package tracker

import (
	"errors"

	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
)

type UUIDGenerator interface {
	Generate() (string, error)
}

var ErrRouteNotFound = errors.New("route not found")

type RouteRepository interface {
	Save(route *domain.Route) error
	Get(uuid domain.RouteUUID) (*domain.Route, error)
}

var ErrActivityNotFound = errors.New("activity not found")

type ActivityRepository interface {
	Save(activity *domain.Activity) error
	Get(uuid domain.ActivityUUID) (*domain.Activity, error)
}

type UserToActivityRepository interface {
	Assign(userUUID auth.UserUUID, activityUUID domain.ActivityUUID) error
	List(userUUID auth.UserUUID) (ActivityIterator, error)
	ListAfter(userUUID auth.UserUUID, startAfter domain.ActivityUUID) (ActivityIterator, error)
	ListBefore(userUUID auth.UserUUID, startBefore domain.ActivityUUID) (ActivityIterator, error)
}

type ActivityWithRoute struct {
	Activity *domain.Activity
	Route    *domain.Route
}

type ActivityIterator interface {
	// Call next in a loop in order to retrieve further items from the
	// iterator.
	Next() (*domain.Activity, bool)

	// After next returns false make sure to call this method to check for
	// errors.
	Error() error
}

type Tracker struct {
	AddActivity        *AddActivityHandler
	GetActivity        *GetActivityHandler
	ListUserActivities *ListUserActivitiesHandler
}

type TransactionProvider interface {
	Read(handler TransactionHandler) error
	Write(handler TransactionHandler) error
}

type TransactionHandler func(repositories *TransactableRepositories) error

type TransactableRepositories struct {
	Route          RouteRepository
	Activity       ActivityRepository
	UserToActivity UserToActivityRepository
}
