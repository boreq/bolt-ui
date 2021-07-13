package tracker

import (
	"fmt"
	"io"
	"strings"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
)

type RouteFileParser interface {
	Parse(r io.Reader, format RouteFileFormat) ([]domain.Point, error)
}

func NewRouteFileFormatFromExtension(extension string) (RouteFileFormat, error) {
	if strings.EqualFold(extension, ".gpx") {
		return RouteFileFormatGpx, nil
	}
	if strings.EqualFold(extension, ".fit") {
		return RouteFileFormatFit, nil
	}
	return RouteFileFormat{}, fmt.Errorf("unknown route file extension '%s'", extension)
}

var (
	RouteFileFormatGpx = RouteFileFormat{"gpx"}
	RouteFileFormatFit = RouteFileFormat{"fit"}
)

type RouteFileFormat struct {
	format string
}

func (f RouteFileFormat) IsZero() bool {
	return f == RouteFileFormat{}
}

func (f RouteFileFormat) String() string {
	return f.format
}

type AddActivity struct {
	RouteFile       io.Reader
	RouteFileFormat RouteFileFormat
	UserUUID        auth.UserUUID
	Visibility      domain.ActivityVisibility
	Title           domain.ActivityTitle
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
	points, err := h.routeFileParser.Parse(cmd.RouteFile, cmd.RouteFileFormat)
	if err != nil {
		return domain.ActivityUUID{}, errors.Wrap(err, "could not parse the route file")
	}

	route, err := h.createRoute(points)
	if err != nil {
		return domain.ActivityUUID{}, errors.Wrap(err, "could not create a route")
	}

	activity, err := h.createActivity(route, cmd.UserUUID, cmd.Visibility, cmd.Title)
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

		if err := repositories.UserToActivity.Assign(activity.UserUUID(), activity.UUID()); err != nil {
			return errors.Wrap(err, "could not map an activity to a user")
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

func (h *AddActivityHandler) createActivity(route *domain.Route, userUUID auth.UserUUID, visibility domain.ActivityVisibility, title domain.ActivityTitle) (*domain.Activity, error) {
	uuid, err := h.uuidGenerator.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate a uuid")
	}

	activityUUID, err := domain.NewActivityUUID(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not create an activity uuid")
	}

	return domain.NewActivity(activityUUID, userUUID, route.UUID(), visibility, title)
}
