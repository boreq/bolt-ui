package domain

import "github.com/pkg/errors"

const maxPrivacyZoneNameLength = 200

type PrivacyZoneName struct {
	s string
}

func NewPrivacyZoneName(s string) (PrivacyZoneName, error) {
	if len(s) > maxPrivacyZoneNameLength {
		return PrivacyZoneName{}, errors.Errorf("length of privacy zone name can not exceed %d characters", maxPrivacyZoneNameLength)
	}

	return PrivacyZoneName{s}, nil
}

func MustNewPrivacyZoneName(s string) PrivacyZoneName {
	t, err := NewPrivacyZoneName(s)
	if err != nil {
		panic(err)
	}
	return t
}

func (n PrivacyZoneName) String() string {
	return n.s
}
