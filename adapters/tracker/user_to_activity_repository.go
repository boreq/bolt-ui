package tracker

import (
	"bytes"

	"github.com/boreq/errors"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/logging"
	bolt "go.etcd.io/bbolt"
)

type UserToActivityRepository struct {
	tx                 *bolt.Tx
	bucket             []byte
	log                logging.Logger
	activityRepository tracker.ActivityRepository
}

func NewUserToActivityRepository(tx *bolt.Tx, activityRepository tracker.ActivityRepository) (*UserToActivityRepository, error) {
	bucket := []byte("user_activities")

	if tx.Writable() {
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return nil, errors.Wrap(err, "could not create a bucket")
		}
	}

	return &UserToActivityRepository{
		tx:                 tx,
		bucket:             bucket,
		log:                logging.New("UserActivityRepository"),
		activityRepository: activityRepository,
	}, nil
}

func (r *UserToActivityRepository) Assign(userUUID auth.UserUUID, activityUUID domain.ActivityUUID) error {
	b, err := r.getOrCreateUserBucket(userUUID)
	if err != nil {
		return errors.Wrap(err, "could not get a bucket")
	}

	return b.Put(activityKey(activityUUID), nil)
}

func (r *UserToActivityRepository) Unassign(userUUID auth.UserUUID, activityUUID domain.ActivityUUID) error {
	b, err := r.getOrCreateUserBucket(userUUID)
	if err != nil {
		return errors.Wrap(err, "could not get a bucket")
	}

	return b.Delete(activityKey(activityUUID))
}

func (r *UserToActivityRepository) List(userUUID auth.UserUUID) (tracker.ActivityIterator, error) {
	b := r.getUserBucket(userUUID)
	if b == nil {
		r.log.Debug("bucket does not exist, returning an empty iterator")
		return newEmptyIterator(), nil
	}

	c := b.Cursor()

	return newActivityIterator(c, r.activityRepository), nil
}

func (r *UserToActivityRepository) ListAfter(userUUID auth.UserUUID, startAfter domain.ActivityUUID) (tracker.ActivityIterator, error) {
	b := r.getUserBucket(userUUID)
	if b == nil {
		r.log.Debug("bucket does not exist, returning an empty iterator")
		return newEmptyIterator(), nil
	}

	c := b.Cursor()

	return newAfterIterator(c, r.activityRepository, startAfter), nil
}

func (r *UserToActivityRepository) ListBefore(userUUID auth.UserUUID, startBefore domain.ActivityUUID) (tracker.ActivityIterator, error) {
	b := r.getUserBucket(userUUID)
	if b == nil {
		r.log.Debug("bucket does not exist, returning an empty iterator")
		return newEmptyIterator(), nil
	}

	c := b.Cursor()

	return newBeforeIterator(c, r.activityRepository, startBefore), nil
}

func activityKey(activityUUID domain.ActivityUUID) []byte {
	return []byte(activityUUID.String())
}

func (r *UserToActivityRepository) getUserBucket(userUUID auth.UserUUID) *bolt.Bucket {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return nil
	}
	return b.Bucket([]byte(userUUID.String()))
}

func (r *UserToActivityRepository) getOrCreateUserBucket(userUUID auth.UserUUID) (*bolt.Bucket, error) {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return nil, errors.New("master bucket does not exist")
	}

	return b.CreateBucketIfNotExists([]byte(userUUID.String()))
}

type activityIterator struct {
	cursor             *bolt.Cursor
	initialized        bool
	activityRepository tracker.ActivityRepository
	err                error
}

func newActivityIterator(cursor *bolt.Cursor, activityRepository tracker.ActivityRepository) *activityIterator {

	return &activityIterator{
		cursor:             cursor,
		activityRepository: activityRepository,
	}

}

func (i *activityIterator) Next() (*domain.Activity, bool) {
	activity, err := i.next()
	if err != nil {
		i.err = err
		return nil, false
	}

	if activity == nil {
		return nil, false
	}

	return activity, true
}

func (i *activityIterator) Error() error {
	return i.err
}

func (i *activityIterator) next() (*domain.Activity, error) {
	var key []byte

	if !i.initialized {
		key, _ = i.cursor.Last()
		i.initialized = true
	} else {
		key, _ = i.cursor.Prev()
	}

	if key == nil {
		return nil, nil
	}

	activityUUID, err := domain.NewActivityUUID(string(key))
	if err != nil {
		return nil, errors.Wrap(err, "could not create a uuid")
	}

	return i.activityRepository.Get(activityUUID)
}

type emptyIterator struct {
}

func newEmptyIterator() *emptyIterator {
	return &emptyIterator{}

}

func (i *emptyIterator) Next() (*domain.Activity, bool) {
	return nil, false
}

func (i *emptyIterator) Error() error {
	return nil
}

type afterIterator struct {
	cursor             *bolt.Cursor
	startAfter         domain.ActivityUUID
	initialized        bool
	activityRepository tracker.ActivityRepository
	err                error
}

func newAfterIterator(cursor *bolt.Cursor, activityRepository tracker.ActivityRepository, startAfter domain.ActivityUUID) *afterIterator {
	return &afterIterator{
		cursor:             cursor,
		activityRepository: activityRepository,
		startAfter:         startAfter,
	}

}

func (i *afterIterator) Next() (*domain.Activity, bool) {
	activity, err := i.next()
	if err != nil {
		i.err = err
		return nil, false
	}

	if activity == nil {
		return nil, false
	}

	return activity, true
}

func (i *afterIterator) Error() error {
	return i.err
}

func (i *afterIterator) next() (*domain.Activity, error) {
	var key []byte

	if !i.initialized {
		searchedKey := activityKey(i.startAfter)
		foundKey, _ := i.cursor.Seek(searchedKey)

		// prevent some weird scanning attacks on the activity list
		if !bytes.Equal(foundKey, searchedKey) {
			return nil, errors.New("unknown activity uuid")
		}

		i.initialized = true
	}

	key, _ = i.cursor.Prev()

	if key == nil {
		return nil, nil
	}

	activityUUID, err := domain.NewActivityUUID(string(key))
	if err != nil {
		return nil, errors.Wrap(err, "could not create a uuid")
	}

	return i.activityRepository.Get(activityUUID)
}

type beforeIterator struct {
	cursor             *bolt.Cursor
	startBefore        domain.ActivityUUID
	initialized        bool
	activityRepository tracker.ActivityRepository
	err                error
}

func newBeforeIterator(cursor *bolt.Cursor, activityRepository tracker.ActivityRepository, startBefore domain.ActivityUUID) *beforeIterator {
	return &beforeIterator{
		cursor:             cursor,
		activityRepository: activityRepository,
		startBefore:        startBefore,
	}

}

func (i *beforeIterator) Next() (*domain.Activity, bool) {
	activity, err := i.next()
	if err != nil {
		i.err = err
		return nil, false
	}

	if activity == nil {
		return nil, false
	}

	return activity, true
}

func (i *beforeIterator) Error() error {
	return i.err
}

func (i *beforeIterator) next() (*domain.Activity, error) {
	var key []byte

	if !i.initialized {
		searchedKey := activityKey(i.startBefore)
		foundKey, _ := i.cursor.Seek(searchedKey)

		// prevent some weird scanning attacks on the activity list
		if !bytes.Equal(foundKey, searchedKey) {
			return nil, errors.New("unknown activity uuid")
		}

		i.initialized = true
	}

	key, _ = i.cursor.Next()

	if key == nil {
		return nil, nil
	}

	activityUUID, err := domain.NewActivityUUID(string(key))
	if err != nil {
		return nil, errors.Wrap(err, "could not create a uuid")
	}

	return i.activityRepository.Get(activityUUID)
}
