package domain

import (
	"fmt"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/internal/eventsourcing"
)

type Activity struct {
	uuid      ActivityUUID
	userUUID  auth.UserUUID
	routeUUID RouteUUID

	es eventsourcing.EventSourcing
}

func NewActivity(uuid ActivityUUID, userUUID auth.UserUUID, routeUUID RouteUUID) (*Activity, error) {
	if uuid.IsZero() {
		return nil, errors.New("zero value of uuid")
	}

	if userUUID.IsZero() {
		return nil, errors.New("zero value of user uuid")
	}

	if routeUUID.IsZero() {
		return nil, errors.New("zero value of route uuid")
	}

	activity := &Activity{}

	if err := activity.update(
		ActivityCreated{
			UUID:      uuid,
			UserUUID:  userUUID,
			RouteUUID: routeUUID,
		},
	); err != nil {
		return nil, errors.Wrap(err, "could not consume the initial event")

	}

	return activity, nil
}

func NewActivityFromHistory(events eventsourcing.EventSourcingEvents) (*Activity, error) {
	activity := &Activity{}

	for _, event := range events {
		if err := activity.update(event.Event); err != nil {
			return nil, errors.Wrap(err, "could not consume an event")
		}
		activity.es.LoadVersion(event)
	}

	activity.PopChanges()

	return activity, nil
}

func (a Activity) UUID() ActivityUUID {
	return a.uuid
}

func (a Activity) RouteUUID() RouteUUID {
	return a.routeUUID
}

func (a Activity) UserUUID() auth.UserUUID {
	return a.userUUID
}

func (a *Activity) PopChanges() eventsourcing.EventSourcingEvents {
	return a.es.PopChanges()
}

func (a *Activity) update(event eventsourcing.Event) error {
	switch e := event.(type) {
	case ActivityCreated:
		a.handleActivityCreated(e)
	default:
		return fmt.Errorf("unknown event type '%T'", event)
	}

	return a.es.Record(event)
}

func (a *Activity) handleActivityCreated(event ActivityCreated) {
	a.uuid = event.UUID
	a.userUUID = event.UserUUID
	a.routeUUID = event.RouteUUID
}

type ActivityCreated struct {
	UUID      ActivityUUID
	UserUUID  auth.UserUUID
	RouteUUID RouteUUID
}

func (e ActivityCreated) EventType() eventsourcing.EventType {
	return "ActivityCreated_v1"
}
