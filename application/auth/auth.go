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
	Generate(username string) (AccessToken, error)
	GetUsername(token AccessToken) (string, error)
}

type PasswordHasher interface {
	Hash(password string) (PasswordHash, error)
	Compare(hashedPassword PasswordHash, password string) error
}

type UserRepository interface {
	// Put inserts the user into the repository. The previous entry with
	// this username is overwriten.
	Put(user User) error

	// Get returns the user with the provided username. If the user doesn't
	// exist ErrNotFound is returned.
	Get(username string) (*User, error)

	// Remove should remove a user.
	Remove(username string) error

	// List should return a list of all users.
	List() ([]User, error)

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

type AccessToken string

type InvitationToken string

type PasswordHash []byte

type User struct {
	UUID          auth.UserUUID
	Username      string
	DisplayName   string
	Password      PasswordHash
	Administrator bool
	Created       time.Time
	LastSeen      time.Time
	Sessions      []Session
}

type Session struct {
	Token    AccessToken
	LastSeen time.Time
}

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
}

func toReadUsers(users []User) []auth.ReadUser {
	var rv []auth.ReadUser
	for _, user := range users {
		rv = append(rv, toReadUser(user))
	}
	return rv
}

func toReadUser(user User) auth.ReadUser {
	rv := auth.ReadUser{
		UUID:          user.UUID,
		Username:      user.Username,
		DisplayName:   user.DisplayName,
		Administrator: user.Administrator,
		Created:       user.Created,
		LastSeen:      user.LastSeen,
	}
	for _, session := range user.Sessions {
		rv.Sessions = append(rv.Sessions, auth.ReadSession{
			LastSeen: session.LastSeen,
		})
	}
	return rv
}
