package domain

import (
	"fmt"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/internal/eventsourcing"
)

type PrivacyZone struct {
	uuid     PrivacyZoneUUID
	userUUID auth.UserUUID
	position Position
	circle   Circle
	name     PrivacyZoneName

	es eventsourcing.EventSourcing
}

func NewPrivacyZone(uuid PrivacyZoneUUID, userUUID auth.UserUUID, position Position, circle Circle, name PrivacyZoneName) (*PrivacyZone, error) {
	if uuid.IsZero() {
		return nil, errors.New("zero value of uuid")
	}

	if userUUID.IsZero() {
		return nil, errors.New("zero value of user uuid")
	}

	if !circle.Contains(position) {
		return nil, errors.New("position is not within the circle")
	}

	zone := &PrivacyZone{}

	if err := zone.update(
		PrivacyZoneCreated{
			UUID:     uuid,
			UserUUID: userUUID,
			Position: position,
			Circle:   circle,
			Name:     name,
		},
	); err != nil {
		return nil, errors.Wrap(err, "could not consume the initial event")

	}

	return zone, nil
}

func NewPrivacyZoneFromHistory(events eventsourcing.EventSourcingEvents) (*PrivacyZone, error) {
	zone := &PrivacyZone{}

	for _, event := range events {
		if err := zone.update(event.Event); err != nil {
			return nil, errors.Wrap(err, "could not consume an event")
		}
		zone.es.LoadVersion(event)
	}

	zone.es.PopChanges()

	return zone, nil
}

func (z *PrivacyZone) UUID() PrivacyZoneUUID {
	return z.uuid
}

func (z *PrivacyZone) Position() Position {
	return z.position
}

func (z *PrivacyZone) Circle() Circle {
	return z.circle
}

func (z *PrivacyZone) UserUUID() auth.UserUUID {
	return z.userUUID
}

func (z *PrivacyZone) Name() PrivacyZoneName {
	return z.name
}

func (z *PrivacyZone) PopChanges() eventsourcing.EventSourcingEvents {
	return z.es.PopChanges()
}

func (z *PrivacyZone) update(event eventsourcing.Event) error {
	switch e := event.(type) {
	case PrivacyZoneCreated:
		z.handlePrivacyZoneCreated(e)
	default:
		return fmt.Errorf("unknown event type '%T'", event)
	}

	return z.es.Record(event)
}

func (z *PrivacyZone) handlePrivacyZoneCreated(event PrivacyZoneCreated) {
	z.uuid = event.UUID
	z.userUUID = event.UserUUID
	z.position = event.Position
	z.circle = event.Circle
	z.name = event.Name
}

type PrivacyZoneCreated struct {
	UUID     PrivacyZoneUUID
	UserUUID auth.UserUUID
	Position Position
	Circle   Circle
	Name     PrivacyZoneName
}

func (e PrivacyZoneCreated) EventType() eventsourcing.EventType {
	return "PrivacyZoneCreated_v1"
}
