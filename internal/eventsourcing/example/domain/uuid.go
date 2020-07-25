package domain

import "errors"

type BankAccountUUID struct {
	uuid string
}

func NewBankAccountUUID(uuid string) (BankAccountUUID, error) {
	if uuid == "" {
		return BankAccountUUID{}, errors.New("uuid can not be empty")
	}

	return BankAccountUUID{uuid}, nil
}

func (u BankAccountUUID) String() string {
	return u.uuid
}

func (u BankAccountUUID) IsZero() bool {
	return u == BankAccountUUID{}
}
