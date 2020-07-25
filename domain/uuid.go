package domain

import "github.com/boreq/errors"

type uuid struct {
	u string
}

func (u uuid) IsZero() bool {
	return u.u == ""
}

func (u uuid) String() string {
	return u.u
}

func newUUID(u string) (uuid, error) {
	if u == "" {
		return uuid{}, errors.New("uuid can not be empty")
	}

	return uuid{u: u}, nil
}

type ActivityUUID struct {
	uuid
}

func NewActivityUUID(u string) (ActivityUUID, error) {
	uuid, err := newUUID(u)
	if err != nil {
		return ActivityUUID{}, errors.New("could not create an activity UUID")
	}
	return ActivityUUID{uuid}, nil
}

func MustNewActivityUUID(u string) ActivityUUID {
	v, err := NewActivityUUID(u)
	if err != nil {
		panic(err)
	}
	return v
}

type UserUUID struct {
	uuid
}

func NewUserUUID(u string) (UserUUID, error) {
	uuid, err := newUUID(u)
	if err != nil {
		return UserUUID{}, errors.New("could not create a user UUID")
	}
	return UserUUID{uuid}, nil
}

func MustNewUserUUID(u string) UserUUID {
	v, err := NewUserUUID(u)
	if err != nil {
		panic(err)
	}
	return v
}

type RouteUUID struct {
	uuid
}

func NewRouteUUID(u string) (RouteUUID, error) {
	uuid, err := newUUID(u)
	if err != nil {
		return RouteUUID{}, errors.New("could not create a route UUID")
	}
	return RouteUUID{uuid}, nil
}

func MustNewRouteUUID(u string) RouteUUID {
	v, err := NewRouteUUID(u)
	if err != nil {
		panic(err)
	}
	return v
}
