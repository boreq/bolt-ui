package domain

import (
	"fmt"

	"github.com/boreq/errors"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/internal/eventsourcing"
)

type Activity struct {
	uuid       ActivityUUID
	userUUID   auth.UserUUID
	routeUUID  RouteUUID
	visibility ActivityVisibility
	title      ActivityTitle

	es eventsourcing.EventSourcing
}

func NewActivity(uuid ActivityUUID, userUUID auth.UserUUID, routeUUID RouteUUID, visibility ActivityVisibility, title ActivityTitle) (*Activity, error) {
	if uuid.IsZero() {
		return nil, errors.New("zero value of uuid")
	}

	if userUUID.IsZero() {
		return nil, errors.New("zero value of user uuid")
	}

	if routeUUID.IsZero() {
		return nil, errors.New("zero value of route uuid")
	}

	if visibility.IsZero() {
		return nil, errors.New("zero value of visibility")
	}

	activity := &Activity{}

	if err := activity.update(
		ActivityCreated{
			UUID:       uuid,
			UserUUID:   userUUID,
			RouteUUID:  routeUUID,
			Visibility: visibility,
			Title:      title,
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

	activity.es.PopChanges()

	return activity, nil
}

func (a *Activity) ChangeTitle(title ActivityTitle) error {
	if title == a.title {
		return nil
	}

	return a.update(
		TitleChanged{
			Title: title,
		},
	)
}

func (a *Activity) ChangeVisibility(visibility ActivityVisibility) error {
	if visibility == a.visibility {
		return nil
	}

	if visibility.IsZero() {
		return errors.New("zero value of visibility")
	}

	return a.update(
		VisibilityChanged{
			Visibility: visibility,
		},
	)
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

func (a Activity) Visibility() ActivityVisibility {
	return a.visibility
}

func (a Activity) Title() ActivityTitle {
	return a.title
}

func (a *Activity) PopChanges() eventsourcing.EventSourcingEvents {
	return a.es.PopChanges()
}

func (a *Activity) HasChanges() bool {
	return a.es.HasChanges()
}

func (a *Activity) update(event eventsourcing.Event) error {
	switch e := event.(type) {
	case ActivityCreated:
		a.handleActivityCreated(e)
	case TitleChanged:
		a.handleTitleChanged(e)
	case VisibilityChanged:
		a.handleVisibilityChanged(e)
	default:
		return fmt.Errorf("unknown event type '%T'", event)
	}

	return a.es.Record(event)
}

func (a *Activity) handleActivityCreated(event ActivityCreated) {
	a.uuid = event.UUID
	a.userUUID = event.UserUUID
	a.routeUUID = event.RouteUUID
	a.visibility = event.Visibility
	a.title = event.Title
}

func (a *Activity) handleTitleChanged(event TitleChanged) {
	a.title = event.Title
}

func (a *Activity) handleVisibilityChanged(event VisibilityChanged) {
	a.visibility = event.Visibility
}

type ActivityCreated struct {
	UUID       ActivityUUID
	UserUUID   auth.UserUUID
	RouteUUID  RouteUUID
	Visibility ActivityVisibility
	Title      ActivityTitle
}

func (e ActivityCreated) EventType() eventsourcing.EventType {
	return "ActivityCreated_v1"
}

type TitleChanged struct {
	Title ActivityTitle
}

func (e TitleChanged) EventType() eventsourcing.EventType {
	return "TitleChanged_v1"
}

type VisibilityChanged struct {
	Visibility ActivityVisibility
}

func (e VisibilityChanged) EventType() eventsourcing.EventType {
	return "VisibilityChanged_v1"
}
