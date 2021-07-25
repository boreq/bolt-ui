package http

import (
	"net/http"

	"github.com/boreq/velo/application"
)

type AuthProvider interface {
	Check(r *http.Request) (bool, error)
}

type TokenAuthProvider struct {
	app *application.Application
}

func NewTokenAuthProvider(app *application.Application) *TokenAuthProvider {
	return &TokenAuthProvider{
		app: app,
	}
}

func (h *TokenAuthProvider) Check(r *http.Request) (bool, error) {
	//token := h.getToken(r)
	//if token == "" {
	//	return nil, nil
	//}

	//cmd := auth.CheckAccessToken{
	//	Token: token,
	//}

	//user, err := h.app.Auth.CheckAccessToken.Execute(cmd)
	//if err != nil {
	//	if errors.Is(err, auth.ErrUnauthorized) {
	//		return nil, nil
	//	}
	//	return nil, errors.Wrap(err, "could not check the access token")
	//}

	//u := AuthenticatedUser{
	//	User:  *user,
	//	Token: token,
	//}

	return true, nil
}

//func (h *TokenAuthProvider) getToken(r *http.Request) authDomain.AccessToken {
//	return authDomain.AccessToken(r.Header.Get("Access-Token"))
//}
