package domain

import "errors"

type Owner struct {
	name string
}

func NewOwner(name string) (Owner, error) {
	if name == "" {
		return Owner{}, errors.New("name can not be empty")
	}

	return Owner{name}, nil
}

func (o Owner) Name() string {
	return o.name
}

func (o Owner) IsZero() bool {
	return o == Owner{}
}
