package tracker

import (
	"errors"

	appAuth "github.com/boreq/velo/application/auth"
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
	Delete(uuid domain.RouteUUID) error
}

var ErrActivityNotFound = errors.New("activity not found")
var ErrGettingActivityForbidden = errors.New("this user can not view this activity")
var ErrEditingActivityForbidden = errors.New("this user can not edit this activity")
var ErrDeletingActivityForbidden = errors.New("this user can not delete this activity")

type ActivityRepository interface {
	Save(activity *domain.Activity) error
	Get(uuid domain.ActivityUUID) (*domain.Activity, error)
	Delete(uuid domain.ActivityUUID) error
}

type UserToActivityRepository interface {
	Assign(userUUID auth.UserUUID, activityUUID domain.ActivityUUID) error
	Unassign(userUUID auth.UserUUID, activityUUID domain.ActivityUUID) error
	List(userUUID auth.UserUUID) (ActivityIterator, error)
	ListAfter(userUUID auth.UserUUID, startAfter domain.ActivityUUID) (ActivityIterator, error)
	ListBefore(userUUID auth.UserUUID, startBefore domain.ActivityUUID) (ActivityIterator, error)
}

type UserRepository interface {
	GetByUUID(uuid auth.UserUUID) (*appAuth.User, error)
}

type Activity struct {
	Activity *domain.Activity
	Route    *domain.Route
	User     *User
}

type User struct {
	Username string
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
	EditActivity       *EditActivityHandler
	DeleteActivity     *DeleteActivityHandler
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
	User           UserRepository
}

func toUser(user *appAuth.User) *User {
	return &User{
		Username: user.Username,
	}
}
