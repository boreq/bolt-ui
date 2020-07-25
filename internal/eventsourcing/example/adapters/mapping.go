package adapters

import (
	"encoding/json"

	"github.com/boreq/velo/internal/eventsourcing"
	"github.com/boreq/velo/internal/eventsourcing/example/domain"
	"github.com/boreq/errors"
)

var mapping = eventsourcing.Mapping{
	"created_v1": eventsourcing.EventMapping{
		Marshal: func(event eventsourcing.Event) ([]byte, error) {
			e := event.(domain.Created)

			transportEvent := created{
				UUID:  e.UUID.String(),
				Owner: e.Owner.Name(),
			}

			return json.Marshal(transportEvent)
		},
		Unmarshal: func(bytes []byte) (eventsourcing.Event, error) {
			var transportEvent created

			if err := json.Unmarshal(bytes, &transportEvent); err != nil {
				return nil, errors.Wrap(err, "could not unmarshal json")
			}

			uuid, err := domain.NewBankAccountUUID(transportEvent.UUID)
			if err != nil {
				return nil, errors.Wrap(err, "could not create uuid")
			}

			owner, err := domain.NewOwner(transportEvent.Owner)
			if err != nil {
				return nil, errors.Wrap(err, "could not create owner")
			}

			return domain.Created{
				UUID:  uuid,
				Owner: owner,
			}, nil
		},
	},
	"deposited_v1": eventsourcing.EventMapping{
		Marshal: func(event eventsourcing.Event) ([]byte, error) {
			e := event.(domain.Deposited)

			transportEvent := deposited{
				Money: e.Money.Amount(),
			}

			return json.Marshal(transportEvent)
		},
		Unmarshal: func(bytes []byte) (eventsourcing.Event, error) {
			var transportEvent deposited

			if err := json.Unmarshal(bytes, &transportEvent); err != nil {
				return nil, errors.Wrap(err, "could not unmarshal json")
			}

			return domain.Deposited{
				Money: domain.NewMoney(transportEvent.Money),
			}, nil
		},
	},
	"withdrawn_v1": eventsourcing.EventMapping{
		Marshal: func(event eventsourcing.Event) ([]byte, error) {
			e := event.(domain.Withdrawn)

			transportEvent := withdrawn{
				Money: e.Money.Amount(),
			}

			return json.Marshal(transportEvent)
		},
		Unmarshal: func(bytes []byte) (eventsourcing.Event, error) {
			var transportEvent withdrawn

			if err := json.Unmarshal(bytes, &transportEvent); err != nil {
				return nil, errors.Wrap(err, "could not unmarshal json")
			}

			return domain.Withdrawn{
				Money: domain.NewMoney(transportEvent.Money),
			}, nil
		},
	},
}

type created struct {
	UUID  string `json:"uuid"`
	Owner string `json:"owner"`
}

type deposited struct {
	Money int `json:"money"`
}

type withdrawn struct {
	Money int `json:"money"`
}
