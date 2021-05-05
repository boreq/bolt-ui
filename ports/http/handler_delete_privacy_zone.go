package http

import (
	"net/http"

	"github.com/boreq/errors"
	"github.com/boreq/rest"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/julienschmidt/httprouter"
)

func (h *Handler) deletePrivacyZone(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())
	uuid := ps.ByName("uuid")

	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	privacyZoneUUID, err := domain.NewPrivacyZoneUUID(uuid)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid privacy zone UUID.")
	}

	cmd := tracker.DeletePrivacyZone{
		PrivacyZoneUUID: privacyZoneUUID,
		AsUser:          &u.User,
	}

	if err := h.app.Tracker.DeletePrivacyZone.Execute(cmd); err != nil {
		if errors.Is(err, tracker.ErrDeletingPrivacyZoneForbidden) {
			return rest.ErrForbidden.WithMessage("You do not have permissions to delete this privacy zone.")
		}

		if errors.Is(err, tracker.ErrPrivacyZoneNotFound) {
			return rest.ErrNotFound.WithMessage("Privacy zone not found.")
		}

		h.log.Error("delete privacy zone command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}
