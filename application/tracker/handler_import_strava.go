package tracker

import (
	"io"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
)

type StravaExportFileParser interface {
	Parse(r io.ReaderAt, size int64) (<-chan StravaActivity, error)
}

type StravaActivity struct {
	Title domain.ActivityTitle
	Route []domain.Point

	Err error
}

type ImportStrava struct {
	Archive     io.ReaderAt
	ArchiveSize int64
	UserUUID    auth.UserUUID
}

type ImportStravaHandler struct {
	transactionProvider    TransactionProvider
	uuidGenerator          UUIDGenerator
	stravaExportFileParser StravaExportFileParser
}

func NewImportStravaHandler(transactionProvider TransactionProvider, uuidGenerator UUIDGenerator, stravaExportFileParser StravaExportFileParser) *ImportStravaHandler {
	return &ImportStravaHandler{
		transactionProvider:    transactionProvider,
		uuidGenerator:          uuidGenerator,
		stravaExportFileParser: stravaExportFileParser,
	}
}

func (h *ImportStravaHandler) Execute(cmd ImportStrava) error {
	ch, err := h.stravaExportFileParser.Parse(cmd.Archive, cmd.ArchiveSize)
	if err != nil {
		return errors.Wrap(err, "failed to parse the strava archive")
	}

	var itemsToSave []activityAndRoute

	for stravaActivity := range ch {
		if stravaActivity.Err != nil {
			return errors.Wrap(stravaActivity.Err, "error iterating over the strava archive")
		}

		route, err := h.createRoute(stravaActivity.Route)
		if err != nil {
			return errors.Wrap(err, "could not create a route")
		}

		activity, err := h.createActivity(route, cmd.UserUUID, stravaActivity.Title)
		if err != nil {
			return errors.Wrap(err, "could not create an activity")
		}

		itemsToSave = append(itemsToSave,
			activityAndRoute{
				Activity: activity,
				Route:    route,
			},
		)
	}

	if err := h.transactionProvider.Write(func(repositories *TransactableRepositories) error {
		for _, item := range itemsToSave {
			if err := repositories.Route.Save(item.Route); err != nil {
				return errors.Wrap(err, "could not save a route")
			}

			if err := repositories.Activity.Save(item.Activity); err != nil {
				return errors.Wrap(err, "could not save an activity ")
			}

			if err := repositories.UserToActivity.Assign(item.Activity.UserUUID(), item.Activity.UUID()); err != nil {
				return errors.Wrap(err, "could not map an activity to a user")
			}
		}

		return nil
	}); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}

func (h *ImportStravaHandler) createRoute(points []domain.Point) (*domain.Route, error) {
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

func (h *ImportStravaHandler) createActivity(route *domain.Route, userUUID auth.UserUUID, title domain.ActivityTitle) (*domain.Activity, error) {
	uuid, err := h.uuidGenerator.Generate()
	if err != nil {
		return nil, errors.Wrap(err, "could not generate a uuid")
	}

	activityUUID, err := domain.NewActivityUUID(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not create an activity uuid")
	}

	return domain.NewActivity(activityUUID, userUUID, route.UUID(), domain.PrivateActivityVisibility, title)
}

type activityAndRoute struct {
	Activity *domain.Activity
	Route    *domain.Route
}
