package http

import (
	"net/http"

	"github.com/boreq/errors"
	"github.com/boreq/velo/application"
	"github.com/boreq/velo/application/auth"
	authDomain "github.com/boreq/velo/domain/auth"
)

type HttpAuthProvider struct {
	app *application.Application
}

func NewHttpAuthProvider(app *application.Application) *HttpAuthProvider {
	return &HttpAuthProvider{
		app: app,
	}
}

func (h *HttpAuthProvider) Get(r *http.Request) (*AuthenticatedUser, error) {
	token := h.getToken(r)
	if token == "" {
		return nil, nil
	}

	cmd := auth.CheckAccessToken{
		Token: token,
	}

	user, err := h.app.Auth.CheckAccessToken.Execute(cmd)
	if err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "could not check the access token")
	}

	u := AuthenticatedUser{
		User:  *user,
		Token: token,
	}

	return &u, nil
}

func (h *HttpAuthProvider) getToken(r *http.Request) authDomain.AccessToken {
	return authDomain.AccessToken(r.Header.Get("Access-Token"))
}
