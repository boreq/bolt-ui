package http

import (
	"encoding/json"
	"net/http"

	"github.com/boreq/errors"
	"github.com/boreq/rest"
	"github.com/boreq/velo/application/auth"
	authDomain "github.com/boreq/velo/domain/auth"
	"github.com/julienschmidt/httprouter"
)

type putUserRequest struct {
	DisplayName string `json:"displayName"`
}

func (h *Handler) putUser(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())
	username := ps.ByName("username")

	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	var t putUserRequest
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.log.Warn("put user decoding failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Malformed input.")
	}

	displayName, err := authDomain.NewDisplayName(t.DisplayName)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid display name.")
	}

	cmd := auth.UpdateProfile{
		Username:    username,
		DisplayName: displayName,
		AsUser:      &u.User,
	}

	if err := h.app.Auth.UpdateProfile.Execute(cmd); err != nil {
		if errors.Is(err, auth.ErrUpdatingProfileForbidden) {
			return rest.ErrForbidden.WithMessage("You do not have permissions to update this user.")
		}

		h.log.Error("update profile command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}
