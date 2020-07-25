package domain

import (
	"fmt"

	"github.com/boreq/velo/internal/eventsourcing"
	"github.com/boreq/errors"
)

type BankAccount struct {
	uuid    BankAccountUUID
	owner   Owner
	balance Money

	es eventsourcing.EventSourcing
}

func NewBankAccount(uuid BankAccountUUID, owner Owner) (*BankAccount, error) {
	if uuid.IsZero() {
		return nil, errors.New("zero value of uuid")
	}

	if owner.IsZero() {
		return nil, errors.New("zero value of owner")
	}

	bankAccount := &BankAccount{}

	event := Created{
		UUID:  uuid,
		Owner: owner,
	}
	if err := bankAccount.update(event); err != nil {
		return nil, errors.Wrap(err, "could not process the event")
	}

	return bankAccount, nil
}

func NewBankAccountFromHistory(events []eventsourcing.EventSourcingEvent) (*BankAccount, error) {
	bankAccount := &BankAccount{}

	for _, event := range events {
		if err := bankAccount.update(event.Event); err != nil {
			return nil, errors.Wrap(err, "could not process the event")
		}
		bankAccount.es.LoadVersion(event)
	}

	bankAccount.es.PopChanges()

	return bankAccount, nil
}

func (a *BankAccount) Deposit(money Money) error {
	event := Deposited{money}
	return a.update(event)
}

func (a *BankAccount) Withdraw(money Money) error {
	if a.balance.Substract(money).IsNegative() {
		return errors.New("you can not withdraw that much money")
	}

	event := Withdrawn{money}
	return a.update(event)
}

func (a *BankAccount) UUID() BankAccountUUID {
	return a.uuid
}

func (a *BankAccount) Owner() Owner {
	return a.owner
}

func (a *BankAccount) Balance() Money {
	return a.balance
}

func (a *BankAccount) HasChanges() bool {
	return a.es.HasChanges()
}

func (a *BankAccount) PopChanges() []eventsourcing.EventSourcingEvent {
	return a.es.PopChanges()
}

func (a *BankAccount) CurrentVersion() eventsourcing.AggregateVersion {
	return a.es.CurrentVersion
}

func (a *BankAccount) update(event eventsourcing.Event) error {
	switch e := event.(type) {
	case Created:
		a.handleCreated(e)
	case Deposited:
		a.handleDeposited(e)
	case Withdrawn:
		a.handleWithdrawn(e)
	default:
		return fmt.Errorf("unsupported event type '%T'", event)
	}

	return a.es.Record(event)
}

func (a *BankAccount) handleCreated(e Created) {
	a.uuid = e.UUID
	a.owner = e.Owner
}

func (a *BankAccount) handleDeposited(e Deposited) {
	a.balance = a.balance.Add(e.Money)
}

func (a *BankAccount) handleWithdrawn(e Withdrawn) {
	a.balance = a.balance.Substract(e.Money)
}

type Created struct {
	UUID  BankAccountUUID
	Owner Owner
}

func (Created) EventType() eventsourcing.EventType {
	return "created_v1"
}

type Deposited struct {
	Money Money
}

func (Deposited) EventType() eventsourcing.EventType {
	return "deposited_v1"
}

type Withdrawn struct {
	Money Money
}

func (Withdrawn) EventType() eventsourcing.EventType {
	return "withdrawn_v1"
}
