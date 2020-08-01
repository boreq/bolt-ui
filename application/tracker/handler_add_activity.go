package tracker

import (
	"io"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
)

type RouteFileParser interface {
	Parse(io.Reader) ([]domain.Point, error)
}

type AddActivity struct {
	RouteFile io.Reader
	UserUUID  domain.UserUUID
}

type AddActivityHandler struct {
	transactionProvider TransactionProvider
	routeFileParser     RouteFileParser
	uuidGenerator       UUIDGenerator
}

func NewAddActivityHandler(transactionProvider TransactionProvider, routeFileParser RouteFileParser, uuidGenerator UUIDGenerator) *AddActivityHandler {
	return &AddActivityHandler{
		transactionProvider: transactionProvider,
		routeFileParser:     routeFileParser,
		uuidGenerator:       uuidGenerator,
	}
}

func (h *AddActivityHandler) Execute(cmd AddActivity) (domain.ActivityUUID, error) {
	points, err := h.routeFileParser.Parse(cmd.RouteFile)
	if err != nil {
		return domain.ActivityUUID{}, errors.Wrap(err, "could not parse the route file")
	}

	route, err := h.createRoute(points)
	if err != nil {
		return domain.ActivityUUID{}, errors.Wrap(err, "could not create a route")
	}

	activity, err := h.createActivity(route, cmd.UserUUID)
	if err != nil {
		return domain.ActivityUUID{}, errors.Wrap(err, "could not create an activity")
	}

	if err := h.transactionProvider.Write(func(repositories *TransactableRepositories) error {
		if err := repositories.Route.Save(route); err != nil {
			return errors.Wrap(err, "could not save a route")
		}

		if err := repositories.Activity.Save(activity); err != nil {
			return errors.Wrap(err, "could not save an activity ")
		}

		return nil
	}); err != nil {
		return domain.ActivityUUID{}, errors.Wrap(err, "transaction failed")
	}

	return activity.UUID(), nil
}

func (h *AddActivityHandler) createRoute(points []domain.Point) (*domain.Route, error) {
	uuid, err := h.uuidGenerator.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate a uuid")
	}

	routeUUID, err := domain.NewRouteUUID(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a route uuid")
	}

	return domain.NewRoute(routeUUID, points)
}

func (h *AddActivityHandler) createActivity(route *domain.Route, userUUID domain.UserUUID) (*domain.Activity, error) {
	uuid, err := h.uuidGenerator.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate a uuid")
	}

	activityUUID, err := domain.NewActivityUUID(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not create an activity uuid")
	}

	return domain.NewActivity(activityUUID, userUUID, route.UUID())
}
