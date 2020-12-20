package auth

import (
	"errors"
	"fmt"
	"regexp"
)

const maxUsernameLen = 100

var usernameRegexp = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

type ValidatedUsername struct {
	username string
}

func NewValidatedUsername(username string) (ValidatedUsername, error) {
	if username == "" {
		return ValidatedUsername{}, errors.New("username can't be empty")
	}

	if len(username) > maxUsernameLen {
		return ValidatedUsername{}, fmt.Errorf("username length can't exceed %d characters", maxUsernameLen)
	}

	if valid := usernameRegexp.MatchString(username); !valid {
		return ValidatedUsername{}, fmt.Errorf("username must conform to the following regexp: %s", usernameRegexp.String())
	}

	return ValidatedUsername{username}, nil
}

func MustNewValidatedUsername(username string) ValidatedUsername {
	v, err := NewValidatedUsername(username)
	if err != nil {
		panic(err)
	}
	return v
}

func (u ValidatedUsername) String() string {
	return u.username
}

func (u ValidatedUsername) IsZero() bool {
	return u == ValidatedUsername{}
}
