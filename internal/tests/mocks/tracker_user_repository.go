package mocks

import (
	"errors"

	authDomain "github.com/boreq/velo/domain/auth"
)

type TrackerUserRepositoryMock struct {
	Users map[authDomain.UserUUID]authDomain.User
}

func NewTrackerUserRepositoryMock() *TrackerUserRepositoryMock {
	return &TrackerUserRepositoryMock{
		Users: make(map[authDomain.UserUUID]authDomain.User),
	}
}

func (r *TrackerUserRepositoryMock) GetByUUID(uuid authDomain.UserUUID) (*authDomain.User, error) {
	u, ok := r.Users[uuid]
	if !ok {
		return nil, errors.New("user not found")
	}

	return &u, nil
}
