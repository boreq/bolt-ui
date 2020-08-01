package auth

import "github.com/boreq/errors"

type uuid struct {
	u string
}

func (u uuid) IsZero() bool {
	return u.u == ""
}

func (u uuid) String() string {
	return u.u
}

func newUUID(u string) (uuid, error) {
	if u == "" {
		return uuid{}, errors.New("uuid can not be empty")
	}

	return uuid{u: u}, nil
}

type UserUUID struct {
	uuid
}

func NewUserUUID(u string) (UserUUID, error) {
	uuid, err := newUUID(u)
	if err != nil {
		return UserUUID{}, errors.New("could not create a user UUID")
	}
	return UserUUID{uuid}, nil
}

func MustNewUserUUID(u string) UserUUID {
	v, err := NewUserUUID(u)
	if err != nil {
		panic(err)
	}
	return v
}
