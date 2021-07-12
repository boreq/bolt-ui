package auth

import (
	"errors"
	"time"
)

type User struct {
	uuid          UserUUID
	username      Username
	displayName   DisplayName
	password      PasswordHash
	administrator bool
	created       time.Time
	lastSeen      time.Time
	sessions      []Session
}

func NewUser(
	uuid UserUUID,
	username Username,
	displayName DisplayName,
	password PasswordHash,
	administrator bool,
) (User, error) {
	if uuid.IsZero() {
		return User{}, errors.New("zero value of uuid")
	}

	if username.IsZero() {
		return User{}, errors.New("zero value of username")
	}

	if displayName.IsZero() {
		return User{}, errors.New("zero value of display name")
	}

	if password.IsZero() {
		return User{}, errors.New("zero value of password")
	}

	return User{
		uuid:          uuid,
		username:      username,
		displayName:   displayName,
		password:      password,
		administrator: administrator,
		created:       time.Now(),
		lastSeen:      time.Now(),
	}, nil
}

func MustNewUser(
	uuid UserUUID,
	username Username,
	displayName DisplayName,
	password PasswordHash,
	administrator bool,
) User {
	u, err := NewUser(
		uuid,
		username,
		displayName,
		password,
		administrator,
	)
	if err != nil {
		panic(err)
	}
	return u
}

func NewHistoricalUser(
	uuid UserUUID,
	username Username,
	displayName DisplayName,
	password PasswordHash,
	administrator bool,
	created time.Time,
	lastSeen time.Time,
	sessions []Session,
) (User, error) {
	if uuid.IsZero() {
		return User{}, errors.New("zero value of uuid")
	}

	if username.IsZero() {
		return User{}, errors.New("zero value of username")
	}

	if displayName.IsZero() {
		return User{}, errors.New("zero value of display name")
	}

	if password.IsZero() {
		return User{}, errors.New("zero value of password")
	}

	if created.IsZero() {
		return User{}, errors.New("zero value of created")
	}

	if lastSeen.IsZero() {
		return User{}, errors.New("zero value of last seen")
	}

	for _, session := range sessions {
		if session.IsZero() {
			return User{}, errors.New("zero value of session")
		}
	}

	return User{
		uuid:          uuid,
		username:      username,
		displayName:   displayName,
		password:      password,
		administrator: administrator,
		created:       created,
		lastSeen:      lastSeen,
		sessions:      copySessions(sessions),
	}, nil
}

func (u User) UUID() UserUUID {
	return u.uuid
}

func (u User) Username() Username {
	return u.username
}

func (u User) DisplayName() DisplayName {
	return u.displayName
}

func (u User) Password() PasswordHash {
	return u.password
}

func (u User) Administrator() bool {
	return u.administrator
}

func (u User) Created() time.Time {
	return u.created
}

func (u User) LastSeen() time.Time {
	return u.lastSeen
}

func (u User) Sessions() []Session {
	return copySessions(u.sessions)
}

func (u *User) ChangeDisplayName(displayName DisplayName) error {
	if displayName.IsZero() {
		return errors.New("zero value of display name")
	}

	u.displayName = displayName

	return nil
}

func (u *User) ChangePassword(password PasswordHash) error {
	if password.IsZero() {
		return errors.New("zero value of password")
	}

	u.password = password

	return nil
}

func (u *User) CheckAccessToken(token AccessToken) (bool, error) {
	if token.IsZero() {
		return false, errors.New("zero value of access token")
	}

	for i := range u.sessions {
		if u.sessions[i].Token() == token {
			u.lastSeen = time.Now()
			u.sessions[i].UpdateLastSeen(time.Now())
			return true, nil
		}
	}

	return false, nil
}

func (u *User) Login(token AccessToken) error {
	if token.IsZero() {
		return errors.New("zero value of access token")
	}

	for i := range u.sessions {
		if u.sessions[i].Token() == token {
			return errors.New("this session already exists")
		}
	}

	s, err := NewSession(token, time.Now())
	if err != nil {
		return errors.New("could not create a session")
	}

	u.sessions = append(u.sessions, s)

	return nil
}

func (u *User) Logout(token AccessToken) error {
	if token.IsZero() {
		return errors.New("zero value of access token")
	}

	for i := range u.sessions {
		if u.sessions[i].Token() == token {
			u.sessions = append(u.sessions[:i], u.sessions[i+1:]...)
			return nil
		}
	}

	return errors.New("session not found")
}

func (u User) AsReadUser() ReadUser {
	rv := ReadUser{
		UUID:          u.uuid,
		Username:      u.username.String(),
		DisplayName:   u.displayName.String(),
		Administrator: u.administrator,
		Created:       u.created,
		LastSeen:      u.lastSeen,
	}
	for _, session := range u.sessions {
		rv.Sessions = append(rv.Sessions, ReadSession{
			LastSeen: session.LastSeen(),
		})
	}
	return rv
}

type Session struct {
	token    AccessToken
	lastSeen time.Time
}

func NewSession(token AccessToken, lastSeen time.Time) (Session, error) {
	if token.IsZero() {
		return Session{}, errors.New("zero value of access token")
	}

	if lastSeen.IsZero() {
		return Session{}, errors.New("zero value of last seen")
	}

	return Session{
		token:    token,
		lastSeen: lastSeen,
	}, nil
}

func (s Session) Token() AccessToken {
	return s.token
}

func (s Session) LastSeen() time.Time {
	return s.lastSeen
}

func (s *Session) UpdateLastSeen(t time.Time) {
	s.lastSeen = t
}

func (s Session) IsZero() bool {
	return s == Session{}
}

type AccessToken string

func (t AccessToken) IsZero() bool {
	return len(t) == 0
}

type PasswordHash []byte

func (h PasswordHash) IsZero() bool {
	return len(h) == 0
}

func copySessions(sessions []Session) []Session {
	tmp := make([]Session, len(sessions))
	copy(tmp, sessions)
	return tmp
}
