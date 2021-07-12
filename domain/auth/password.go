package auth

import (
	"errors"
	"fmt"
)

const maxPasswordLen = 10000

type Password struct {
	password string
}

func NewPassword(password string) (Password, error) {
	if password == "" {
		return Password{}, errors.New("password can't be empty")
	}

	if len(password) > maxPasswordLen {
		return Password{}, fmt.Errorf("password length can't exceed %d characters", maxPasswordLen)
	}

	return Password{password}, nil
}

func MustNewPassword(password string) Password {
	v, err := NewPassword(password)
	if err != nil {
		panic(err)
	}
	return v
}

func (u Password) String() string {
	return u.password
}

func (u Password) IsZero() bool {
	return u == Password{}
}
