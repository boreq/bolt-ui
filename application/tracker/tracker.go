package tracker

import (
	"errors"

	"github.com/boreq/velo/domain"
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

type Tracker struct {
	AddActivity *AddActivityHandler
	GetActivity *GetActivityHandler
}

type TransactionProvider interface {
	Read(handler TransactionHandler) error
	Write(handler TransactionHandler) error
}

type TransactionHandler func(repositories *TransactableRepositories) error

type TransactableRepositories struct {
	Route    RouteRepository
	Activity ActivityRepository
}
