package example

import (
	"testing"

	"github.com/boreq/eggplant/internal/eventsourcing/example/adapters"
	"github.com/boreq/eggplant/internal/eventsourcing/example/domain"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	uuid, err := domain.NewBankAccountUUID("bank-account-uuid")
	require.NoError(t, err)

	owner, err := domain.NewOwner("bank-account-owner")
	require.NoError(t, err)

	bankAccount, err := domain.NewBankAccount(uuid, owner)
	require.NoError(t, err)

	err = bankAccount.Deposit(domain.NewMoney(100))
	require.NoError(t, err)

	err = bankAccount.Withdraw(domain.NewMoney(10))
	require.NoError(t, err)

	err = bankAccount.Deposit(domain.NewMoney(20))
	require.NoError(t, err)

	balance := bankAccount.Balance()
	require.Equal(t, domain.NewMoney(110), balance)

	repository := adapters.NewBankAccountRepository()

	err = repository.Save(bankAccount)
	require.NoError(t, err)

	loadedBankAccount, err := repository.Get(uuid)
	require.NoError(t, err)

	require.False(t, loadedBankAccount.HasChanges())

	require.Equal(t, bankAccount.UUID(), loadedBankAccount.UUID())
	require.Equal(t, bankAccount.Owner(), loadedBankAccount.Owner())
	require.Equal(t, bankAccount.Balance(), loadedBankAccount.Balance())
	require.Equal(t, bankAccount.CurrentVersion(), loadedBankAccount.CurrentVersion())
}
