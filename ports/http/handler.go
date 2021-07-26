package http

import (
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/boreq/errors"
	"github.com/boreq/rest"
	"github.com/boreq/velo/application"
	"github.com/boreq/velo/logging"
	"github.com/boreq/velo/ports/http/frontend"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	app          *application.Application
	authProvider AuthProvider
	router       *httprouter.Router
	log          logging.Logger
}

func NewHandler(app *application.Application, authProvider AuthProvider) (*Handler, error) {
	h := &Handler{
		app:          app,
		authProvider: authProvider,
		router:       httprouter.New(),
		log:          logging.New("ports/http.Handler"),
	}

	h.router.HandlerFunc(http.MethodGet, "/api/browse/*path", rest.Wrap(h.browse))

	ffs, err := frontend.NewFrontendFileSystem()
	if err != nil {
		return nil, err
	}
	h.router.NotFound = http.FileServer(ffs)

	return h, nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(rw, req)
}

func (h *Handler) browse(r *http.Request) rest.RestResponse {
	ps := httprouter.ParamsFromContext(r.Context())

	ok, err := h.authProvider.Check(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if !ok {
		return rest.ErrForbidden.WithMessage("Invalid token.")
	}

	path, err := readPath(ps.ByName("path"))
	if err != nil {
		h.log.Warn("invalid path", "err", err)
		return rest.ErrBadRequest.WithMessage("Invalid path.")
	}

	query := application.Browse{
		Path: path,
	}

	tree, err := h.app.Browse.Execute(query)
	if err != nil {
		if errors.Is(err, application.ErrBucketNotFound) {
			return rest.ErrNotFound
		}
		h.log.Error("browse failure", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(
		toTree(tree),
	)
}

const sep = "/"

func readPath(s string) ([]application.Key, error) {
	s = strings.Trim(s, sep)

	if s == "" {
		return nil, nil
	}

	var path []application.Key

	for _, element := range strings.Split(s, "/") {
		b, err := hex.DecodeString(element)
		if err != nil {
			return nil, errors.Wrap(err, "could not decode")
		}

		key, err := application.NewKey(b)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a key")
		}

		path = append(path, key)
	}

	return path, nil
}
