package http

import (
	"encoding/json"
	"net/http"

	"github.com/boreq/rest"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
)

type PostPrivacyZoneRequest struct {
	Position Position `json:"position"`
	Circle   Circle   `json:"circle"`
	Name     string   `json:"name"`
}

type PostPrivacyZoneResponse struct {
	PrivacyZoneUUID string `json:"privacyZoneUUID"`
}

func (h *Handler) postPrivacyZone(r *http.Request) rest.RestResponse {
	log := h.handlerLogger("post privacy zone")

	u, err := h.authProvider.Get(r)
	if err != nil {
		log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	var j PostPrivacyZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		log.Warn("decoding failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Malformed input.")
	}

	position, err := fromPosition(j.Position)
	if err != nil {
		log.Warn("invalid position", "err", err)
		return rest.ErrBadRequest.WithMessage("Invalid position.")
	}

	circle, err := fromCircle(j.Circle)
	if err != nil {
		log.Warn("invalid circle", "err", err)
		return rest.ErrBadRequest.WithMessage("Invalid circle.")
	}

	name, err := domain.NewPrivacyZoneName(j.Name)
	if err != nil {
		log.Warn("invalid name", "err", err)
		return rest.ErrBadRequest.WithMessage("Invalid name.")
	}

	cmd := tracker.AddPrivacyZone{
		UserUUID: u.User.UUID,
		Position: position,
		Circle:   circle,
		Name:     name,
	}

	privacyZoneUUID, err := h.app.Tracker.AddPrivacyZone.Execute(cmd)
	if err != nil {
		log.Error("add privacy zone command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(
		PostPrivacyZoneResponse{
			PrivacyZoneUUID: privacyZoneUUID.String(),
		},
	)
}
