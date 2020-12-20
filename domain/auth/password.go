package auth

import (
	"errors"
	"fmt"
)

const maxPasswordLen = 10000

type ValidatedPassword struct {
	password string
}

func NewValidatedPassword(password string) (ValidatedPassword, error) {
	if password == "" {
		return ValidatedPassword{}, errors.New("password can't be empty")
	}

	if len(password) > maxPasswordLen {
		return ValidatedPassword{}, fmt.Errorf("password length can't exceed %d characters", maxPasswordLen)
	}

	return ValidatedPassword{password}, nil
}

func MustNewValidatedPassword(password string) ValidatedPassword {
	v, err := NewValidatedPassword(password)
	if err != nil {
		panic(err)
	}
	return v
}

func (u ValidatedPassword) String() string {
	return u.password
}

func (u ValidatedPassword) IsZero() bool {
	return u == ValidatedPassword{}
}
