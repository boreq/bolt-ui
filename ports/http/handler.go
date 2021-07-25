package http

import (
	"net/http"

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
	//	ps := httprouter.ParamsFromContext(r.Context())
	//	path := strings.Trim(ps.ByName("path"), "/")

	ok, err := h.authProvider.Check(r)
	if err != nil {
		h.log.Error("auth provider get failed", "err", err)
		return rest.ErrInternalServerError
	}

	if !ok {
		return rest.ErrForbidden.WithMessage("Invalid token.")
	}

	//query := application.Browse{}

	//if err := h.app.Browse.Execute(query); err != nil {
	//	h.log.Error("could not browse", "err", err)
	//	return rest.ErrInternalServerError
	//}

	return rest.NewResponse(nil)
}
