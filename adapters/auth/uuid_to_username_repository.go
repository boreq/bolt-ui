package auth

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/application/auth"
	authDomain "github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/logging"
	bolt "go.etcd.io/bbolt"
)

type UUIDToUsernameRepository struct {
	tx     *bolt.Tx
	bucket []byte
	log    logging.Logger
}

func NewUUIDToUsernameRepository(tx *bolt.Tx) (*UUIDToUsernameRepository, error) {
	bucket := []byte("users_usernames")

	if tx.Writable() {
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return nil, errors.Wrap(err, "could not create a bucket")
		}
	}

	return &UUIDToUsernameRepository{
		tx:     tx,
		bucket: bucket,
		log:    logging.New("UUIDToUsernameRepository"),
	}, nil
}

func (r *UUIDToUsernameRepository) Remove(uuid authDomain.UserUUID) error {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return errors.New("bucket does not exist")
	}
	return b.Delete([]byte(uuid.String()))
}

func (r *UUIDToUsernameRepository) Get(uuid authDomain.UserUUID) (string, error) {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return "", errors.Wrap(auth.ErrNotFound, "bucket does not exist")
	}
	u := b.Get([]byte(uuid.String()))
	if u == nil {
		return "", auth.ErrNotFound
	}

	return string(u), nil
}

func (r *UUIDToUsernameRepository) Put(uuid authDomain.UserUUID, username authDomain.Username) error {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return errors.New("bucket does not exist")
	}
	return b.Put([]byte(uuid.String()), []byte(username.String()))
}
