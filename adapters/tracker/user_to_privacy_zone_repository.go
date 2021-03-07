package tracker

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/logging"
	bolt "go.etcd.io/bbolt"
)

type UserToPrivacyZoneRepository struct {
	tx                    *bolt.Tx
	bucket                []byte
	log                   logging.Logger
	privacyZoneRepository tracker.PrivacyZoneRepository
}

func NewUserToPrivacyZoneRepository(tx *bolt.Tx, privacyZoneRepository tracker.PrivacyZoneRepository) (*UserToPrivacyZoneRepository, error) {
	bucket := []byte("user_privacy_zones")

	if tx.Writable() {
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return nil, errors.Wrap(err, "could not create a bucket")
		}
	}

	return &UserToPrivacyZoneRepository{
		tx:                    tx,
		bucket:                bucket,
		log:                   logging.New("UserToPrivacyZoneRepository"),
		privacyZoneRepository: privacyZoneRepository,
	}, nil
}

func (r *UserToPrivacyZoneRepository) Assign(userUUID auth.UserUUID, privacyZoneUUID domain.PrivacyZoneUUID) error {
	b, err := r.getOrCreateUserBucket(userUUID)
	if err != nil {
		return errors.Wrap(err, "could not get a bucket")
	}

	return b.Put(r.privacyZoneKey(privacyZoneUUID), nil)
}

func (r *UserToPrivacyZoneRepository) Unassign(userUUID auth.UserUUID, privacyZoneUUID domain.PrivacyZoneUUID) error {
	b, err := r.getOrCreateUserBucket(userUUID)
	if err != nil {
		return errors.Wrap(err, "could not get a bucket")
	}

	return b.Delete(r.privacyZoneKey(privacyZoneUUID))
}

func (r *UserToPrivacyZoneRepository) List(userUUID auth.UserUUID) ([]*domain.PrivacyZone, error) {
	return nil, errors.New("not implemented")
}

func (r *UserToPrivacyZoneRepository) privacyZoneKey(uuid domain.PrivacyZoneUUID) []byte {
	return []byte(uuid.String())
}

func (r *UserToPrivacyZoneRepository) getUserBucket(userUUID auth.UserUUID) *bolt.Bucket {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return nil
	}
	return b.Bucket([]byte(userUUID.String()))
}

func (r *UserToPrivacyZoneRepository) getOrCreateUserBucket(userUUID auth.UserUUID) (*bolt.Bucket, error) {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return nil, errors.New("master bucket does not exist")
	}

	return b.CreateBucketIfNotExists([]byte(userUUID.String()))
}
