package http

import (
	"net/http"

	"github.com/boreq/errors"
	"github.com/boreq/rest"
	"github.com/boreq/velo/application/auth"
	"github.com/boreq/velo/application/tracker"
	"github.com/julienschmidt/httprouter"
)

func (h *Handler) getUserPrivacyZones(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())
	username := ps.ByName("username")

	currentUser, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	user, err := h.app.Auth.GetUser.Execute(auth.GetUser{
		Username: username,
	})
	if err != nil {
		h.log.Error("could not get a user", "err", err)
		return rest.ErrInternalServerError
	}

	query := tracker.ListUserPrivacyZones{
		UserUUID: user.UUID,
		AsUser:   currentUser.UserPointer(),
	}

	privacyZones, err := h.app.Tracker.ListUserPrivacyZones.Execute(query)
	if err != nil {
		if errors.Is(err, tracker.ErrGettingPrivacyZoneForbidden) {
			return rest.ErrForbidden.WithMessage("You do not have permissions to list this user's privacy zones.")
		}

		h.log.Error("could not get user activities", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(toPrivacyZones(privacyZones))
}
