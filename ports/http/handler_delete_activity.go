package http

import (
	"net/http"

	"github.com/boreq/errors"
	"github.com/boreq/rest"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
	"github.com/julienschmidt/httprouter"
)

func (h *Handler) deleteActivity(r *http.Request) rest.RestResponse {
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

	activityUUID, err := domain.NewActivityUUID(uuid)
	if err != nil {
		return rest.ErrBadRequest.WithMessage("Invalid activity UUID.")
	}

	cmd := tracker.DeleteActivity{
		ActivityUUID: activityUUID,
		AsUser:       &u.User,
	}

	if err := h.app.Tracker.DeleteActivity.Execute(cmd); err != nil {
		if errors.Is(err, tracker.ErrDeletingActivityForbidden) {
			return rest.ErrForbidden.WithMessage("You do not have permissions to delete this activity.")
		}

		if errors.Is(err, tracker.ErrActivityNotFound) {
			return rest.ErrNotFound.WithMessage("Activity not found.")
		}

		h.log.Error("delete activity command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}
