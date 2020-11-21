package mocks

import (
	"errors"

	"github.com/boreq/velo/application/auth"
	authDomain "github.com/boreq/velo/domain/auth"
)

type TrackerUserRepositoryMock struct {
	Users map[authDomain.UserUUID]auth.User
}

func NewTrackerUserRepositoryMock() *TrackerUserRepositoryMock {
	return &TrackerUserRepositoryMock{
		Users: make(map[authDomain.UserUUID]auth.User),
	}
}

func (r *TrackerUserRepositoryMock) GetByUUID(uuid authDomain.UserUUID) (*auth.User, error) {
	u, ok := r.Users[uuid]
	if !ok {
		return nil, errors.New("user not found")
	}

	return &u, nil
}
