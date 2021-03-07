package tracker

import (
	"encoding/json"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/internal/eventsourcing"
)

var privacyZoneEventMapping = eventsourcing.Mapping{
	"PrivacyZoneCreated_v1": eventsourcing.EventMapping{
		Marshal: func(event eventsourcing.Event) ([]byte, error) {
			e := event.(domain.PrivacyZoneCreated)

			transportEvent := privacyZoneCreated{
				UUID:     e.UUID.String(),
				UserUUID: e.UserUUID.String(),
				Position: toPosition(e.Position),
				Circle:   toCircle(e.Circle),
				Name:     e.Name.String(),
			}

			return json.Marshal(transportEvent)
		},
		Unmarshal: func(bytes []byte) (eventsourcing.Event, error) {
			var transportEvent privacyZoneCreated

			if err := json.Unmarshal(bytes, &transportEvent); err != nil {
				return nil, errors.Wrap(err, "could not unmarshal json")
			}

			uuid, err := domain.NewPrivacyZoneUUID(transportEvent.UUID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create a uuid")
			}

			userUUID, err := auth.NewUserUUID(transportEvent.UserUUID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create a user uuid")
			}

			position, err := fromPosition(transportEvent.Position)
			if err != nil {
				return nil, errors.Wrap(err, "could not map position")
			}

			circle, err := fromCircle(transportEvent.Circle)
			if err != nil {
				return nil, errors.Wrap(err, "could not map circle")
			}

			name, err := domain.NewPrivacyZoneName(transportEvent.Name)
			if err != nil {
				return nil, errors.Wrap(err, "could not create a name")
			}

			return domain.PrivacyZoneCreated{
				UUID:     uuid,
				UserUUID: userUUID,
				Position: position,
				Circle:   circle,
				Name:     name,
			}, nil
		},
	},
}

type privacyZoneCreated struct {
	UUID     string   `json:"uuid"`
	UserUUID string   `json:"userUUID"`
	Position position `json:"position"`
	Circle   circle   `json:"circle"`
	Name     string   `json:"name"`
}
