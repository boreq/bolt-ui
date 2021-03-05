package domain

import "github.com/pkg/errors"

const maxActivityTitleLength = 50

type ActivityTitle struct {
	s string
}

func NewActivityTitle(s string) (ActivityTitle, error) {
	if len(s) > maxActivityTitleLength {
		return ActivityTitle{}, errors.Errorf("length of activity title can not exceed %d characters", maxActivityTitleLength)
	}

	return ActivityTitle{s}, nil
}

func MustNewActivityTitle(s string) ActivityTitle {
	t, err := NewActivityTitle(s)
	if err != nil {
		panic(err)
	}
	return t
}

func (a ActivityTitle) String() string {
	return a.s
}
