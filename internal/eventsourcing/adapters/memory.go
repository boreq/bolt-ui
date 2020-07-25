package adapters

import "github.com/boreq/velo/internal/eventsourcing"

type MemoryPersistenceAdapter struct {
	events map[eventsourcing.AggregateUUID][]eventsourcing.PersistedEvent
}

func NewMemoryPersistenceAdapter() *MemoryPersistenceAdapter {
	return &MemoryPersistenceAdapter{
		events: make(map[eventsourcing.AggregateUUID][]eventsourcing.PersistedEvent),
	}
}

func (m *MemoryPersistenceAdapter) SaveEvents(aggregateUUID eventsourcing.AggregateUUID, events []eventsourcing.PersistedEvent) error {
	if len(events) == 0 {
		return eventsourcing.EmptyEventsErr
	}

	m.events[aggregateUUID] = append(m.events[aggregateUUID], events...)
	return nil
}

func (m *MemoryPersistenceAdapter) GetEvents(aggregateUUID eventsourcing.AggregateUUID) ([]eventsourcing.PersistedEvent, error) {
	if len(m.events[aggregateUUID]) == 0 {
		return nil, eventsourcing.EventsNotFound
	}
	return m.events[aggregateUUID], nil
}
