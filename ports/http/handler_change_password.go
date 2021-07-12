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

type putUserPasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (h *Handler) putUserPassword(r *http.Request) rest.RestResponse {
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

	var t putUserPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.log.Warn("put user password decoding failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Malformed input.")
	}

	newPassword, err := authDomain.NewValidatedPassword(t.NewPassword)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid password.")
	}

	cmd := auth.ChangePassword{
		Username:    username,
		OldPassword: t.OldPassword,
		NewPassword: newPassword,
		AsUser:      &u.User,
	}

	if err := h.app.Auth.ChangePassword.Execute(cmd); err != nil {
		if errors.Is(err, auth.ErrUnauthorized) {
			return rest.ErrForbidden.WithMessage("Invalid old password.")
		}

		if errors.Is(err, auth.ErrChangingPasswordForbidden) {
			return rest.ErrForbidden.WithMessage("You do not have permissions to change this password.")
		}

		h.log.Error("change password command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}
