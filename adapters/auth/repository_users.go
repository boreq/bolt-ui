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
	uuidToUsernameRepository *UUIDToUsernameRepository
	tx                       *bolt.Tx
	bucket                   []byte
	log                      logging.Logger
}

func NewUserRepository(tx *bolt.Tx, uuidToUsernameRepository *UUIDToUsernameRepository) (*UserRepository, error) {
	bucket := []byte("users")

	if tx.Writable() {
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return nil, errors.Wrap(err, "could not create a bucket")
		}
	}

	return &UserRepository{
		uuidToUsernameRepository: uuidToUsernameRepository,
		tx:                       tx,
		bucket:                   bucket,
		log:                      logging.New("UserRepository"),
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

func (r *UserRepository) List() ([]authDomain.User, error) {
	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return nil, nil
	}

	c := b.Cursor()

	var users []authDomain.User
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var pu persistedUser
		if err := json.Unmarshal(v, &pu); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}

		u, err := r.fromPersisted(pu)
		if err != nil {
			return nil, errors.Wrap(err, "could not convert from persisted")
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) Remove(username string) error {
	u, err := r.Get(username)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return nil
		}
		return errors.Wrap(err, "could not get a user")
	}

	if err := r.uuidToUsernameRepository.Remove(u.UUID()); err != nil {
		return errors.Wrap(err, "could not remove the username mapping")
	}

	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return errors.New("bucket does not exist")
	}
	return b.Delete([]byte(username))
}

func (r *UserRepository) Get(username string) (*authDomain.User, error) {
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

	tmp, err := r.fromPersisted(u)
	return &tmp, err
}

func (r *UserRepository) GetByUUID(uuid authDomain.UserUUID) (*authDomain.User, error) {
	username, err := r.uuidToUsernameRepository.Get(uuid)
	if err != nil {
		return nil, errors.Wrap(err, "could not get the mapping")
	}

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

	tmp, err := r.fromPersisted(u)
	return &tmp, err
}

func (r *UserRepository) Put(user authDomain.User) error {
	persistedUser := r.toPersisted(user)

	j, err := json.Marshal(persistedUser)
	if err != nil {
		return errors.Wrap(err, "marshaling to json failed")
	}

	if err := r.uuidToUsernameRepository.Put(user.UUID(), user.Username()); err != nil {
		return errors.Wrap(err, "could not store the mapping")
	}

	b := r.tx.Bucket(r.bucket)
	if b == nil {
		return errors.New("bucket does not exist")
	}
	return b.Put([]byte(user.Username().String()), j)
}

func (r *UserRepository) toPersisted(user authDomain.User) persistedUser {
	var sessions []persistedSession

	for _, session := range user.Sessions() {
		sessions = append(sessions, persistedSession{
			Token:    string(session.Token()),
			LastSeen: session.LastSeen(),
		})
	}

	return persistedUser{
		UUID:          user.UUID().String(),
		Username:      user.Username().String(),
		DisplayName:   user.DisplayName().String(),
		Password:      user.Password(),
		Administrator: user.Administrator(),
		Created:       user.Created(),
		LastSeen:      user.LastSeen(),
		Sessions:      sessions,
	}
}

func (r *UserRepository) fromPersisted(user persistedUser) (authDomain.User, error) {
	uuid, err := authDomain.NewUserUUID(user.UUID)
	if err != nil {
		return authDomain.User{}, errors.Wrap(err, "could not create user uuid")
	}

	username, err := authDomain.NewUsername(user.Username)
	if err != nil {
		return authDomain.User{}, errors.Wrap(err, "could not create username")
	}

	displayName, err := authDomain.NewDisplayName(user.DisplayName)
	if err != nil {
		return authDomain.User{}, errors.Wrap(err, "could not create display name")
	}

	var sessions []authDomain.Session

	for _, session := range user.Sessions {
		session, err := authDomain.NewSession(
			authDomain.AccessToken(session.Token),
			session.LastSeen,
		)
		if err != nil {
			return authDomain.User{}, errors.Wrap(err, "could not create a session")
		}

		sessions = append(sessions, session)
	}

	return authDomain.NewHistoricalUser(
		uuid,
		username,
		displayName,
		authDomain.PasswordHash(user.Password),
		user.Administrator,
		user.Created,
		user.LastSeen,
		sessions,
	)
}
