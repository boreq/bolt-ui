package auth

import (
	"errors"
	"fmt"
)

const maxDisplayNameLen = 100

type DisplayName struct {
	displayName string
}

func NewDisplayName(displayName string) (DisplayName, error) {
	if displayName == "" {
		return DisplayName{}, errors.New("display name can't be empty")
	}

	if len(displayName) > maxDisplayNameLen {
		return DisplayName{}, fmt.Errorf("display name length can't exceed %d characters", maxDisplayNameLen)
	}

	return DisplayName{displayName}, nil
}

func MustNewDisplayName(displayName string) DisplayName {
	v, err := NewDisplayName(displayName)
	if err != nil {
		panic(err)
	}
	return v
}

func (u DisplayName) String() string {
	return u.displayName
}

func (u DisplayName) IsZero() bool {
	return u == DisplayName{}
}
