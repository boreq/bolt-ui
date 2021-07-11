package http

import (
	"net/http"

	"github.com/boreq/rest"
	"github.com/boreq/velo/application/tracker"
)

const maxStravaExportFileSize = 10 * 1024 * 1024 // max size of the strava export file in bytes

func (h *Handler) importStrava(r *http.Request) rest.RestResponse {
	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	file, header, err := r.FormFile("archive")
	if err != nil {
		h.log.Warn("export file retrieval failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Failed to retrieve the file.")
	}

	if header.Size > maxStravaExportFileSize {
		return rest.ErrBadRequest.WithMessage("Activity file too large.")
	}

	cmd := tracker.ImportStrava{
		UserUUID:    u.User.UUID,
		Archive:     file,
		ArchiveSize: header.Size,
	}

	err = h.app.Tracker.ImportStrava.Execute(cmd)
	if err != nil {
		h.log.Error("add activity command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(nil)
}
