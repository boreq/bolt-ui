package auth

import (
	"github.com/boreq/velo/domain/auth"
	"github.com/pkg/errors"
)

type RegisterInitial struct {
	Username auth.Username
	Password auth.Password
}

type RegisterInitialHandler struct {
	passwordHasher      PasswordHasher
	transactionProvider TransactionProvider
	uuidGenerator       UUIDGenerator
}

func NewRegisterInitialHandler(
	passwordHasher PasswordHasher,
	transactionProvider TransactionProvider,
	uuidGenerator UUIDGenerator,
) *RegisterInitialHandler {
	return &RegisterInitialHandler{
		passwordHasher:      passwordHasher,
		transactionProvider: transactionProvider,
		uuidGenerator:       uuidGenerator,
	}
}

func (h *RegisterInitialHandler) Execute(cmd RegisterInitial) error {
	if cmd.Username.IsZero() {
		return errors.New("zero value of username")
	}

	if cmd.Password.IsZero() {
		return errors.New("zero value of password")
	}

	passwordHash, err := h.passwordHasher.Hash(cmd.Password.String())
	if err != nil {
		return errors.Wrap(err, "hashing the password failed")
	}

	uuid, err := h.uuidGenerator.Generate()
	if err != nil {
		return errors.Wrap(err, "could not generate an uuid")
	}

	userUUID, err := auth.NewUserUUID(uuid)
	if err != nil {
		return errors.Wrap(err, "could not create a user uuid")
	}

	displayName, err := auth.NewDisplayName(cmd.Username.String())
	if err != nil {
		return errors.Wrap(err, "could not create a display name")
	}

	u, err := auth.NewUser(
		userUUID,
		cmd.Username,
		displayName,
		passwordHash,
		true,
	)
	if err != nil {
		return errors.Wrap(err, "could not create a user")
	}

	if err := h.transactionProvider.Write(func(r *TransactableRepositories) error {
		n, err := r.Users.Count()
		if err != nil {
			return errors.Wrap(err, "could not get a count")
		}
		if n != 0 {
			return errors.New("there are existing users")
		}
		return r.Users.Put(u)
	}); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}
