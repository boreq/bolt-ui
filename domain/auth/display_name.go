package auth

import (
	"errors"
	"fmt"
)

const maxDisplayNameLen = 100

type ValidatedDisplayName struct {
	displayName string
}

func NewValidatedDisplayName(displayName string) (ValidatedDisplayName, error) {
	if displayName == "" {
		return ValidatedDisplayName{}, errors.New("display name can't be empty")
	}

	if len(displayName) > maxDisplayNameLen {
		return ValidatedDisplayName{}, fmt.Errorf("display name length can't exceed %d characters", maxDisplayNameLen)
	}

	return ValidatedDisplayName{displayName}, nil
}

func MustNewValidatedDisplayName(displayName string) ValidatedDisplayName {
	v, err := NewValidatedDisplayName(displayName)
	if err != nil {
		panic(err)
	}
	return v
}

func (u ValidatedDisplayName) String() string {
	return u.displayName
}

func (u ValidatedDisplayName) IsZero() bool {
	return u == ValidatedDisplayName{}
}
