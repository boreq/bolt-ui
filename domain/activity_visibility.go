package domain

import "fmt"

var (
	PublicActivityVisibility   = ActivityVisibility{"public"}
	UnlistedActivityVisibility = ActivityVisibility{"unlisted"}
	PrivateActivityVisibility  = ActivityVisibility{"private"}
)

var stringToActivityVisibility = map[string]ActivityVisibility{
	PublicActivityVisibility.String():   PublicActivityVisibility,
	UnlistedActivityVisibility.String(): UnlistedActivityVisibility,
	PrivateActivityVisibility.String():  PrivateActivityVisibility,
}

type ActivityVisibility struct {
	s string
}

func NewActivityVisibility(s string) (ActivityVisibility, error) {
	v, ok := stringToActivityVisibility[s]
	if !ok {
		return ActivityVisibility{}, fmt.Errorf("invalid input: %s", s)
	}

	return v, nil
}

func (a ActivityVisibility) String() string {
	return a.s
}

func (a ActivityVisibility) IsZero() bool {
	return a == ActivityVisibility{}
}
