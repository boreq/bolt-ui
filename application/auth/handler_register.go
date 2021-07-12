package auth

import (
	"github.com/boreq/errors"
	"github.com/boreq/velo/domain/auth"
)

type Register struct {
	Username auth.Username
	Password auth.Password
	Token    InvitationToken
}

type RegisterHandler struct {
	passwordHasher      PasswordHasher
	transactionProvider TransactionProvider
	uuidGenerator       UUIDGenerator
}

func NewRegisterHandler(
	passwordHasher PasswordHasher,
	transactionProvider TransactionProvider,
	uuidGenerator UUIDGenerator,
) *RegisterHandler {
	return &RegisterHandler{
		passwordHasher:      passwordHasher,
		transactionProvider: transactionProvider,
		uuidGenerator:       uuidGenerator,
	}
}

func (h *RegisterHandler) Execute(cmd Register) error {
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
		false,
	)
	if err != nil {
		return errors.Wrap(err, "could not create a user")
	}

	if err := h.transactionProvider.Write(func(r *TransactableRepositories) error {
		if _, err := r.Invitations.Get(cmd.Token); err != nil {
			return errors.Wrap(err, "could not get the invitation")
		}

		if err := r.Invitations.Remove(cmd.Token); err != nil {
			return errors.Wrap(err, "could not remove the invitation")
		}

		if _, err := r.Users.Get(cmd.Username.String()); err != nil {
			if !errors.Is(err, ErrNotFound) {
				return errors.Wrap(err, "could not get the user")
			}
		} else {
			return ErrUsernameTaken
		}

		return r.Users.Put(u)
	}); err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}
