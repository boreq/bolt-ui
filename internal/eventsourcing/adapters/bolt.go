package adapters

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boreq/eggplant/internal/eventsourcing"
	"github.com/boreq/errors"
	bolt "go.etcd.io/bbolt"
)

type BucketName []byte
type GetBucketPathFunc func(aggregate eventsourcing.AggregateUUID) []BucketName

type BoltPersistenceAdapter struct {
	tx       *bolt.Tx
	pathFunc GetBucketPathFunc
}

func NewBoltPersistenceAdapter(tx *bolt.Tx, pathFunc GetBucketPathFunc) *BoltPersistenceAdapter {
	return &BoltPersistenceAdapter{
		tx:       tx,
		pathFunc: pathFunc,
	}
}

func (b *BoltPersistenceAdapter) SaveEvents(uuid eventsourcing.AggregateUUID, events []eventsourcing.PersistedEvent) error {
	if len(events) == 0 {
		return eventsourcing.EmptyEventsErr
	}

	bucket, err := b.createAggregateBucket(uuid)
	if err != nil {
		return errors.Wrap(err, "could not get a bucket")
	}

	if err := b.validateEvents(bucket, events); err != nil {
		return errors.Wrap(err, "invalid events")
	}

	for _, event := range events {
		key := toKey(event.AggregateVersion)

		value, err := json.Marshal(event)
		if err != nil {
			return errors.Wrap(err, "event marshaling failed")
		}

		if err := bucket.Put(key, value); err != nil {
			return errors.Wrap(err, "could not call put")
		}
	}

	return nil
}

func (b *BoltPersistenceAdapter) GetEvents(uuid eventsourcing.AggregateUUID) ([]eventsourcing.PersistedEvent, error) {
	bucket, err := b.getAggregateBucket(uuid)
	if err != nil {
		if errors.Is(err, eventsourcing.EventsNotFound) {
			return nil, err
		}
		return nil, errors.Wrap(err, "could not get a bucket")
	}

	var persistedEvents []eventsourcing.PersistedEvent

	c := bucket.Cursor()

	for key, value := c.First(); key != nil; key, value = c.Next() {
		var persistedEvent eventsourcing.PersistedEvent

		if err := json.Unmarshal(value, &persistedEvent); err != nil {
			return nil, errors.Wrap(err, "event unmarshaling failed")
		}

		persistedEvents = append(persistedEvents, persistedEvent)
	}

	if len(persistedEvents) == 0 {
		return nil, eventsourcing.EventsNotFound
	}

	return persistedEvents, nil
}

func (b *BoltPersistenceAdapter) createAggregateBucket(uuid eventsourcing.AggregateUUID) (bucket *bolt.Bucket, err error) {
	bucketNames := b.pathFunc(uuid)
	if len(bucketNames) == 0 {
		return nil, errors.New("path func returned an empty slice")
	}

	return b.createBucket(bucketNames)
}

func (b *BoltPersistenceAdapter) getAggregateBucket(uuid eventsourcing.AggregateUUID) (*bolt.Bucket, error) {
	bucketNames := b.pathFunc(uuid)
	if len(bucketNames) == 0 {
		return nil, errors.New("path func returned an empty slice")
	}

	bucket := b.getBucket(bucketNames)
	if bucket == nil {
		return nil, eventsourcing.EventsNotFound
	}

	return bucket, nil
}

func (b *BoltPersistenceAdapter) validateEvents(bucket *bolt.Bucket, events []eventsourcing.PersistedEvent) error {
	// todo check if event n + 1 has the version of the event n incremented by one

	// todo check if the newest persisted event has the version of the oldest new event decremented by one

	//if len(events) == 0 {
	//	return eventsourcing.EmptyEventsErr
	//}
	//
	//key, value := bucket.Cursor().Last()
	//if key != nil && value != nil {
	//	lastEvent :=
	//
	//}

	return nil
}

func (b *BoltPersistenceAdapter) getBucket(bucketNames []BucketName) *bolt.Bucket {
	bucket := b.tx.Bucket(bucketNames[0])

	if bucket == nil {
		return nil
	}

	for i := 1; i < len(bucketNames); i++ {
		bucket = bucket.Bucket(bucketNames[i])
		if bucket == nil {
			return nil
		}
	}

	return bucket
}

func (b *BoltPersistenceAdapter) createBucket(bucketNames []BucketName) (bucket *bolt.Bucket, err error) {
	bucket, err = b.tx.CreateBucketIfNotExists(bucketNames[0])
	if err != nil {
		return nil, errors.Wrap(err, "could not create a bucket")
	}

	for i := 1; i < len(bucketNames); i++ {
		bucket, err = bucket.CreateBucketIfNotExists(bucketNames[i])
		if err != nil {
			return nil, errors.Wrap(err, "could not create a bucket")
		}
	}

	return bucket, nil
}

func toKey(version eventsourcing.AggregateVersion) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(version))
	return b
}

//func fromKey(key []byte) eventsourcing.AggregateVersion {
//	return eventsourcing.AggregateVersion(binary.BigEndian.Uint64(key))
//}
