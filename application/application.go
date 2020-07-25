package application

import (
	"github.com/boreq/velo/application/auth"
	"github.com/boreq/velo/application/tracker"
)

type Application struct {
	Auth    auth.Auth
	Tracker tracker.Tracker
}
