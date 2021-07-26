package http

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/boreq/velo/internal/config"
	"github.com/boreq/velo/logging"
	"github.com/rs/cors"
)

type Server struct {
	handler http.Handler
	conf    *config.Config
	log     logging.Logger
}

func NewServer(handler http.Handler, conf *config.Config) *Server {
	return &Server{
		handler: handler,
		conf:    conf,
		log:     logging.New("ports/http.Server"),
	}
}

func (s *Server) Serve(address string) error {
	handler := s.handler

	if s.conf.InsecureCORS {
		handler = cors.AllowAll().Handler(s.handler)
	}

	handler = gziphandler.GzipHandler(handler)

	s.log.Debug("starting listening", "address", address)
	return http.ListenAndServe(address, handler)
}
