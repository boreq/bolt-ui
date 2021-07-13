package http

import (
	"net/http"
	"path"

	"github.com/boreq/rest"
	"github.com/boreq/velo/application/tracker"
	"github.com/boreq/velo/domain"
)

const maxActivityFileSize = 10 * 1024 * 1024 // max size of the activity file in bytes

func (h *Handler) postActivity(r *http.Request) rest.RestResponse {
	u, err := h.authProvider.Get(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if u == nil {
		return rest.ErrUnauthorized
	}

	file, header, err := r.FormFile("routeFile")
	if err != nil {
		h.log.Warn("activity file retrieval failed", "err", err)
		return rest.ErrBadRequest.WithMessage("Failed to retrieve the file.")
	}

	if header.Size > maxActivityFileSize {
		return rest.ErrBadRequest.WithMessage("Activity file too large.")
	}

	format, err := tracker.NewRouteFileFormatFromExtension(path.Ext(header.Filename))
	if err != nil {
		h.log.Warn("invalid extension", "err", err)
		return rest.ErrBadRequest.WithMessage("Unsupported file extension.")
	}

	visibility, err := domain.NewActivityVisibility(r.FormValue("visibility"))
	if err != nil {
		h.log.Warn("invalid visiblity", "err", err)
		return rest.ErrBadRequest.WithMessage("Invalid visiblity.")
	}

	title, err := domain.NewActivityTitle(r.FormValue("title"))
	if err != nil {
		h.log.Warn("invalid title", "err", err)
		return rest.ErrBadRequest.WithMessage("Invalid title.")
	}

	cmd := tracker.AddActivity{
		UserUUID:        u.User.UUID,
		RouteFile:       file,
		RouteFileFormat: format,
		Visibility:      visibility,
		Title:           title,
	}

	activityUUID, err := h.app.Tracker.AddActivity.Execute(cmd)
	if err != nil {
		h.log.Error("add activity command failed", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(
		PostActivityResponse{
			ActivityUUID: activityUUID.String(),
		},
	)
}
