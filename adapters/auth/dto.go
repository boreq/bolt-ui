package auth

import "time"

type persistedUser struct {
	UUID          string             `json:"uuid"`
	Username      string             `json:"username"`
	DisplayName   string             `json:"displayName"`
	Password      []byte             `json:"password"`
	Administrator bool               `json:"administrator"`
	Created       time.Time          `json:"created"`
	LastSeen      time.Time          `json:"lastSeen"`
	Sessions      []persistedSession `json:"sessions"`
}

type persistedSession struct {
	Token    string    `json:"token"`
	LastSeen time.Time `json:"lastSeen"`
}

type persistedInvitation struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
}
