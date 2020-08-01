package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/internal/eventsourcing"
	"github.com/boreq/velo/internal/eventsourcing/adapters"
	bolt "go.etcd.io/bbolt"
)

const activityBucket = "activities"

type ActivityRepository struct {
	tx         *bolt.Tx
	eventStore *eventsourcing.EventStore
}

func NewActivityRepository(tx *bolt.Tx) (*ActivityRepository, error) {
	persistenceAdapter := adapters.NewBoltPersistenceAdapter(
		tx,
		func(uuid eventsourcing.AggregateUUID) []adapters.BucketName {
			return []adapters.BucketName{
				[]byte(activityBucket),
				[]byte(uuid),
				[]byte("events"),
			}
		},
	)
	eventStore := eventsourcing.NewEventStore(activityEventMapping, persistenceAdapter)

	return &ActivityRepository{
		tx:         tx,
		eventStore: eventStore,
	}, nil
}

func (c *ActivityRepository) Save(activity *domain.Activity) error {
	return c.eventStore.SaveEvents(c.convertUUID(activity.UUID()), activity.PopChanges())
}

func (c *ActivityRepository) Get(uuid domain.ActivityUUID) (*domain.Activity, error) {
	events, err := c.eventStore.GetEvents(c.convertUUID(uuid))
	if err != nil {
		if errors.Is(err, eventsourcing.EventsNotFound) {
			return nil, tracker.ErrActivityNotFound
		}
		return nil, errors.Wrap(err, "could not get events")
	}

	return domain.NewActivityFromHistory(events)
}

func (c *ActivityRepository) convertUUID(uuid domain.ActivityUUID) eventsourcing.AggregateUUID {
	return eventsourcing.AggregateUUID(uuid.String())
}
