package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/internal/eventsourcing"
	"github.com/boreq/velo/internal/eventsourcing/adapters"
	bolt "go.etcd.io/bbolt"
)

const privacyZoneBucket = "privacy_zones"

type PrivacyZoneRepository struct {
	tx         *bolt.Tx
	eventStore *eventsourcing.EventStore
}

func NewPrivacyZoneRepository(tx *bolt.Tx) (*PrivacyZoneRepository, error) {
	persistenceAdapter := adapters.NewBoltPersistenceAdapter(
		tx,
		func(uuid eventsourcing.AggregateUUID) []adapters.BucketName {
			return []adapters.BucketName{
				[]byte(privacyZoneBucket),
				[]byte(uuid),
				[]byte("events"),
			}
		},
	)
	eventStore := eventsourcing.NewEventStore(privacyZoneEventMapping, persistenceAdapter)

	return &PrivacyZoneRepository{
		tx:         tx,
		eventStore: eventStore,
	}, nil
}

func (c *PrivacyZoneRepository) Save(zone *domain.PrivacyZone) error {
	return c.eventStore.SaveEvents(c.convertUUID(zone.UUID()), zone.PopChanges())
}

func (c *PrivacyZoneRepository) Get(uuid domain.PrivacyZoneUUID) (*domain.PrivacyZone, error) {
	events, err := c.eventStore.GetEvents(c.convertUUID(uuid))
	if err != nil {
		if errors.Is(err, eventsourcing.EventsNotFound) {
			return nil, tracker.ErrPrivacyZoneNotFound
		}
		return nil, errors.Wrap(err, "could not get events")
	}

	return domain.NewPrivacyZoneFromHistory(events)
}

func (c *PrivacyZoneRepository) Delete(uuid domain.PrivacyZoneUUID) error {
	b := c.tx.Bucket([]byte(privacyZoneBucket))
	if b == nil {
		return nil
	}

	return b.DeleteBucket([]byte(uuid.String()))
}

func (c *PrivacyZoneRepository) convertUUID(uuid domain.PrivacyZoneUUID) eventsourcing.AggregateUUID {
	return eventsourcing.AggregateUUID(uuid.String())
}
