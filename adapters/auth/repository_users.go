package auth

import (
	"encoding/json"

	"github.com/boreq/errors"
	"github.com/boreq/velo/application/auth"
	authDomain "github.com/boreq/velo/domain/auth"
	"github.com/boreq/velo/logging"
	bolt "go.etcd.io/bbolt"
)

type UserRepository struct {
	tx     *bolt.Tx
	bucket []byte
	log    logging.Logger
}

func NewUserRepository(tx *bolt.Tx) (*UserRepository, error) {
	bucket := []byte("users")

	if tx.Writable() {
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return nil, errors.Wrap(err, "could not create a bucket")
		}
	}

	return &UserRepository{
		tx:     tx,
		bucket: bucket,
		log:    logging.New("UserRepository"),
	}, nil
}

func (r *UserRepository) Count() (int, error) {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return 0, nil
	}
	count := b.Stats().KeyN
	return count, nil
}

func (r *UserRepository) List() ([]auth.User, error) {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return nil, nil
	}

	c := b.Cursor()

	var users []auth.User
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var pu persistedUser
		if err := json.Unmarshal(v, &pu); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}

		u, err := r.fromPersisted(pu)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert from persisted")
		}

		users = append(users, *u)
	}

	return users, nil
}

func (r *UserRepository) Remove(username string) error {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return errors.New("bucket does not exist")
	}
	return b.Delete([]byte(username))
}

func (r *UserRepository) Get(username string) (*auth.User, error) {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return nil, errors.Wrap(auth.ErrNotFound, "bucket does not exist")
	}
	j := b.Get([]byte(username))
	if j == nil {
		return nil, auth.ErrNotFound
	}

	u := persistedUser{}
	if err := json.Unmarshal(j, &u); err != nil {
		return nil, errors.Wrap(err, "json unmarshal failed")
	}

	return r.fromPersisted(u)
}

func (r *UserRepository) Put(user auth.User) error {
	persistedUser := r.toPersisted(user)

	j, err := json.Marshal(persistedUser)
	if err != nil {
		return errors.Wrap(err, "marshaling to json failed")
	}

	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return errors.New("bucket does not exist")
	}
	return b.Put([]byte(user.Username), j)
}

func (r *UserRepository) toPersisted(user auth.User) persistedUser {
	var sessions []persistedSession

	for _, session := range user.Sessions {
		sessions = append(sessions, persistedSession{
			Token:    string(session.Token),
			LastSeen: session.LastSeen,
		})
	}

	return persistedUser{
		UUID:          user.UUID.String(),
		Username:      user.Username,
		Password:      user.Password,
		Administrator: user.Administrator,
		Created:       user.Created,
		LastSeen:      user.LastSeen,
		Sessions:      sessions,
	}
}

func (r *UserRepository) fromPersisted(user persistedUser) (*auth.User, error) {
	uuid, err := authDomain.NewUserUUID(user.UUID)
	if err != nil {
		return nil, errors.Wrap(err, "could not create user uuid")
	}

	var sessions []auth.Session

	for _, session := range user.Sessions {
		sessions = append(sessions, auth.Session{
			Token:    auth.AccessToken(session.Token),
			LastSeen: session.LastSeen,
		})
	}

	return &auth.User{
		UUID:          uuid,
		Username:      user.Username,
		Password:      user.Password,
		Administrator: user.Administrator,
		Created:       user.Created,
		LastSeen:      user.LastSeen,
		Sessions:      sessions,
	}, nil
}
