package auth

import (
	"time"
)

type ReadUser struct {
	UUID          UserUUID
	Username      string
	Administrator bool
	Created       time.Time
	LastSeen      time.Time
	Sessions      []ReadSession
}

type ReadSession struct {
	LastSeen time.Time
}
