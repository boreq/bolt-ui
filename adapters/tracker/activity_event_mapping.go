package tracker

import (
	"encoding/json"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/internal/eventsourcing"
)

var activityEventMapping = eventsourcing.Mapping{
	"ActivityCreated_v1": eventsourcing.EventMapping{
		Marshal: func(event eventsourcing.Event) ([]byte, error) {
			e := event.(domain.ActivityCreated)

			transportEvent := activityCreated{
				UUID:       e.UUID.String(),
				UserUUID:   e.UserUUID.String(),
				RouteUUID:  e.RouteUUID.String(),
				Visibility: e.Visibility.String(),
				Title:      e.Title.String(),
			}

			return json.Marshal(transportEvent)
		},
		Unmarshal: func(bytes []byte) (eventsourcing.Event, error) {
			var transportEvent activityCreated

			if err := json.Unmarshal(bytes, &transportEvent); err != nil {
				return nil, errors.Wrap(err, "could not unmarshal json")
			}

			uuid, err := domain.NewActivityUUID(transportEvent.UUID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create a uuid")
			}

			userUUID, err := auth.NewUserUUID(transportEvent.UserUUID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create a user uuid")
			}

			routeUUID, err := domain.NewRouteUUID(transportEvent.RouteUUID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create a route uuid")
			}

			visibility, err := domain.NewActivityVisibility(transportEvent.Visibility)
			if err != nil {
				return nil, errors.Wrap(err, "could not create activity visibility")
			}

			title, err := domain.NewActivityTitle(transportEvent.Title)
			if err != nil {
				return nil, errors.Wrap(err, "could not create activity title")
			}

			return domain.ActivityCreated{
				UUID:       uuid,
				UserUUID:   userUUID,
				RouteUUID:  routeUUID,
				Visibility: visibility,
				Title:      title,
			}, nil
		},
	},
}

type activityCreated struct {
	UUID       string `json:"uuid"`
	UserUUID   string `json:"userUUID"`
	RouteUUID  string `json:"routeUUID"`
	Visibility string `json:"visibility"`
	Title      string `json:"title"`
}
