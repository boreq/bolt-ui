package adapters_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/boreq/eggplant/internal/eventsourcing"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/require"
)

type Test func(t *testing.T, adapter eventsourcing.PersistenceAdapter)

type TestRunner func(t *testing.T, test Test)

func TestAdapters(t *testing.T) {
	adapters := []struct {
		Name       string
		TestRunner TestRunner
	}{
		{
			Name:       "memory",
			TestRunner: RunTestMemory,
		},
		{
			Name:       "bolt",
			TestRunner: RunTestBolt,
		},
	}

	tests := []struct {
		Name string
		Test Test
	}{
		{
			Name: "save_empty_events",
			Test: testSaveEmptyEvents,
		},
		{
			Name: "test_save_events",
			Test: testSaveEvents,
		},
	}

	for _, adapter := range adapters {
		t.Run(adapter.Name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.Name, func(t *testing.T) {
					adapter.TestRunner(t, test.Test)
				})
			}
		})
	}
}

func testSaveEmptyEvents(t *testing.T, adapter eventsourcing.PersistenceAdapter) {
	uuid := someAggregateUUID()

	err := adapter.SaveEvents(uuid, nil)
	require.Equal(t, eventsourcing.EmptyEventsErr, err)
}

func testSaveEvents(t *testing.T, adapter eventsourcing.PersistenceAdapter) {
	aggregateUUID := someAggregateUUID()

	events := []eventsourcing.PersistedEvent{
		{
			EventPayload:     nil,
			EventType:        "some_event",
			AggregateVersion: 0,
		},
		{
			EventPayload:     nil,
			EventType:        "some_event",
			AggregateVersion: 1,
		},
	}

	persistedEvents, err := adapter.GetEvents(aggregateUUID)
	require.Equal(t, eventsourcing.EventsNotFound, err)
	require.Empty(t, persistedEvents)

	err = adapter.SaveEvents(aggregateUUID, events)
	require.NoError(t, err)

	persistedEvents, err = adapter.GetEvents(aggregateUUID)
	require.NoError(t, err)
	require.Equal(t, events, persistedEvents)
}

func someAggregateUUID() eventsourcing.AggregateUUID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	ulid := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	return eventsourcing.AggregateUUID(ulid.String())
}
