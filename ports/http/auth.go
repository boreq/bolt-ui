package http

import (
	"errors"
	"net/http"

	"github.com/boreq/velo/internal/config"
)

type AuthProvider interface {
	Check(r *http.Request) (bool, error)
}

type TokenAuthProvider struct {
	conf *config.Config
}

func NewTokenAuthProvider(conf *config.Config) *TokenAuthProvider {
	return &TokenAuthProvider{
		conf: conf,
	}
}

func (h *TokenAuthProvider) Check(r *http.Request) (bool, error) {
	if h.conf.InsecureToken {
		return true, nil
	}

	if h.conf.Token == "" {
		return false, errors.New("auth token is not set in the config")
	}

	token := r.Header.Get("Access-Token")
	if token != h.conf.Token {
		return false, nil
	}

	return true, nil
}
