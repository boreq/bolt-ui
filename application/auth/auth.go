package auth

import (
	"errors"
	"time"

	"github.com/boreq/velo/domain/auth"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrUsernameTaken = errors.New("username taken")
var ErrNotFound = errors.New("not found")

type UUIDGenerator interface {
	Generate() (string, error)
}

type CryptoStringGenerator interface {
	Generate(bytes int) (string, error)
}

type AccessTokenGenerator interface {
	Generate(username string) (auth.AccessToken, error)
	GetUsername(token auth.AccessToken) (string, error)
}

type PasswordHasher interface {
	Hash(password string) (auth.PasswordHash, error)
	Compare(hashedPassword auth.PasswordHash, password string) error
}

type UserRepository interface {
	// Put inserts the user into the repository. The previous entry with
	// this username is overwriten.
	Put(user auth.User) error

	// Get returns the user with the provided username. If the user doesn't
	// exist ErrNotFound is returned.
	Get(username string) (*auth.User, error)

	// Remove should remove a user.
	Remove(username string) error

	// List should return a list of all users.
	List() ([]auth.User, error)

	// Count returns the number of users.
	Count() (int, error)
}

type InvitationRepository interface {
	// Put inserts the invitation into the repository. The previous entry
	// with this token is overwriten.
	Put(invitation Invitation) error

	// Get returns an invitation with the provided token, if the invitation
	// doesn't exist ErrNotFound is returned.
	Get(token InvitationToken) (*Invitation, error)

	// Remove removes an invitation. If the invitation doesn't exist this
	// function returns nil.
	Remove(token InvitationToken) error
}

type InvitationToken string

type Invitation struct {
	Token   InvitationToken
	Created time.Time
}

type TransactionProvider interface {
	Read(handler TransactionHandler) error
	Write(handler TransactionHandler) error
}

type TransactionHandler func(repositories *TransactableRepositories) error

type TransactableRepositories struct {
	Invitations InvitationRepository
	Users       UserRepository
}

type Auth struct {
	RegisterInitial  *RegisterInitialHandler
	Register         *RegisterHandler
	Login            *LoginHandler
	Logout           *LogoutHandler
	CheckAccessToken *CheckAccessTokenHandler
	List             *ListHandler
	CreateInvitation *CreateInvitationHandler
	Remove           *RemoveHandler
	SetPassword      *SetPasswordHandler
	GetUser          *GetUserHandler
	UpdateProfile    *UpdateProfileHandler
	ChangePassword   *ChangePasswordHandler
}

func toReadUsers(users []auth.User) []auth.ReadUser {
	var rv []auth.ReadUser
	for _, user := range users {
		rv = append(rv, user.AsReadUser())
	}
	return rv
}
