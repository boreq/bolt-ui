package auth

import (
	"errors"
	"fmt"
	"regexp"
)

const maxUsernameLen = 100

var usernameRegexp = regexp.MustCompile("^[a-zA-Z0-9_-]+$")

type Username struct {
	username string
}

func NewUsername(username string) (Username, error) {
	if username == "" {
		return Username{}, errors.New("username can't be empty")
	}

	if len(username) > maxUsernameLen {
		return Username{}, fmt.Errorf("username length can't exceed %d characters", maxUsernameLen)
	}

	if valid := usernameRegexp.MatchString(username); !valid {
		return Username{}, fmt.Errorf("username must conform to the following regexp: %s", usernameRegexp.String())
	}

	return Username{username}, nil
}

func MustNewUsername(username string) Username {
	v, err := NewUsername(username)
	if err != nil {
		panic(err)
	}
	return v
}

func (u Username) String() string {
	return u.username
}

func (u Username) IsZero() bool {
	return u == Username{}
}
