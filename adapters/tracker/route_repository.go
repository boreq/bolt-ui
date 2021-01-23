package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/internal/eventsourcing"
	"github.com/boreq/velo/internal/eventsourcing/adapters"
	bolt "go.etcd.io/bbolt"
)

const routeBucket = "routes"

type RouteRepository struct {
	tx         *bolt.Tx
	eventStore *eventsourcing.EventStore
}

func NewRouteRepository(tx *bolt.Tx) (*RouteRepository, error) {
	persistenceAdapter := adapters.NewBoltPersistenceAdapter(
		tx,
		func(uuid eventsourcing.AggregateUUID) []adapters.BucketName {
			return []adapters.BucketName{
				[]byte(routeBucket),
				[]byte(uuid),
				[]byte("events"),
			}
		},
	)
	eventStore := eventsourcing.NewEventStore(routeEventMapping, persistenceAdapter)

	return &RouteRepository{
		tx:         tx,
		eventStore: eventStore,
	}, nil
}

func (c *RouteRepository) Save(route *domain.Route) error {
	return c.eventStore.SaveEvents(c.convertUUID(route.UUID()), route.PopChanges())
}

func (c *RouteRepository) Get(uuid domain.RouteUUID) (*domain.Route, error) {
	events, err := c.eventStore.GetEvents(c.convertUUID(uuid))
	if err != nil {
		if errors.Is(err, eventsourcing.EventsNotFound) {
			return nil, tracker.ErrRouteNotFound
		}
		return nil, errors.Wrap(err, "could not get events")
	}

	return domain.NewRouteFromHistory(events)
}

func (c *RouteRepository) Delete(uuid domain.RouteUUID) error {
	b := c.tx.Bucket([]byte(routeBucket))
	if b == nil {
		return nil
	}

	return b.DeleteBucket([]byte(uuid.String()))
}

func (c *RouteRepository) convertUUID(uuid domain.RouteUUID) eventsourcing.AggregateUUID {
	return eventsourcing.AggregateUUID(uuid.String())
}
