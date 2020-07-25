package eventsourcing

import (
	"fmt"

	"github.com/boreq/errors"
)

var EmptyEventsErr = errors.New("passed empty events slice")
var EventsNotFound = errors.New("events not found")

type EventType string

type AggregateVersion uint64

type AggregateUUID string

type Event interface {
	EventType() EventType
}

type EventSourcingEvent struct {
	Event            Event
	AggregateVersion AggregateVersion
}

type EventSourcingEvents []EventSourcingEvent

func (e EventSourcingEvents) Payloads() []Event {
	var payloads []Event

	for _, event := range e {
		payloads = append(payloads, event.Event)
	}

	return payloads
}

type PersistedEvent struct {
	EventPayload     []byte           `json:"eventPayload"`
	EventType        EventType        `json:"eventType"`
	AggregateVersion AggregateVersion `json:"aggregateVersion"`
}

type PersistenceAdapter interface {
	SaveEvents(aggregateUUID AggregateUUID, events []PersistedEvent) error
	GetEvents(aggregateUUID AggregateUUID) ([]PersistedEvent, error)
}

type Mapping map[EventType]EventMapping

type EventMapping struct {
	Marshal   func(Event) ([]byte, error)
	Unmarshal func([]byte) (Event, error)
}

type EventStore struct {
	mapping            Mapping
	persistenceAdapter PersistenceAdapter
}

func NewEventStore(mapping Mapping, persistenceAdapter PersistenceAdapter) *EventStore {
	return &EventStore{
		mapping:            mapping,
		persistenceAdapter: persistenceAdapter,
	}
}

func (s *EventStore) SaveEvents(aggregateUUID AggregateUUID, events []EventSourcingEvent) error {
	if len(events) == 0 {
		return EmptyEventsErr
	}

	marshaledEvents, err := s.marshalEvents(events)
	if err != nil {
		return errors.Wrap(err, "could not marshal events")
	}

	if err := s.persistenceAdapter.SaveEvents(aggregateUUID, marshaledEvents); err != nil {
		return errors.Wrap(err, "could not save events using the persistence adapter")
	}

	return nil
}

func (s *EventStore) GetEvents(aggregateUUID AggregateUUID) ([]EventSourcingEvent, error) {
	marshaledEvents, err := s.persistenceAdapter.GetEvents(aggregateUUID)
	if err != nil {
		return nil, errors.Wrap(err, "could not get events using the persistence adapter")
	}

	events, err := s.unmarshalEvents(marshaledEvents)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal events")
	}

	return events, nil
}

func (s *EventStore) marshalEvents(events []EventSourcingEvent) ([]PersistedEvent, error) {
	var persistedEvents []PersistedEvent

	for _, event := range events {
		persisitedEvent, err := s.marshalEvent(event)
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal an event")
		}

		persistedEvents = append(persistedEvents, persisitedEvent)
	}

	return persistedEvents, nil
}

func (s *EventStore) marshalEvent(event EventSourcingEvent) (PersistedEvent, error) {
	mapping, ok := s.mapping[event.Event.EventType()]
	if !ok {
		return PersistedEvent{}, fmt.Errorf("missing event mapping for '%s'", event.Event.EventType())
	}

	payload, err := mapping.Marshal(event.Event)
	if err != nil {
		return PersistedEvent{}, errors.Wrapf(err, "could not marshal '%s'", event.Event.EventType())
	}

	return PersistedEvent{
		EventPayload:     payload,
		EventType:        event.Event.EventType(),
		AggregateVersion: event.AggregateVersion,
	}, nil
}

func (s *EventStore) unmarshalEvents(persistedEvents []PersistedEvent) ([]EventSourcingEvent, error) {
	var events []EventSourcingEvent

	for _, persistedEvent := range persistedEvents {
		event, err := s.unmarshalEvent(persistedEvent)
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal an event")
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *EventStore) unmarshalEvent(persistedEvent PersistedEvent) (EventSourcingEvent, error) {
	mapping, ok := s.mapping[persistedEvent.EventType]
	if !ok {
		return EventSourcingEvent{}, fmt.Errorf("missing event mapping for '%s'", persistedEvent.EventType)
	}

	event, err := mapping.Unmarshal(persistedEvent.EventPayload)
	if err != nil {
		return EventSourcingEvent{}, errors.Wrapf(err, "could not unmarshal '%s'", persistedEvent.EventType)
	}

	return EventSourcingEvent{
		Event:            event,
		AggregateVersion: persistedEvent.AggregateVersion,
	}, nil
}
