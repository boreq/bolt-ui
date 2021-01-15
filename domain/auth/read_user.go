package auth

import (
	"time"
)

type ReadUser struct {
	UUID          UserUUID      `json:"-"`
	Username      string        `json:"username"`
	Administrator bool          `json:"administrator"`
	Created       time.Time     `json:"created"`
	LastSeen      time.Time     `json:"lastSeen"`
	Sessions      []ReadSession `json:"sessions"`
}

type ReadSession struct {
	LastSeen time.Time `json:"lastSeen"`
}
