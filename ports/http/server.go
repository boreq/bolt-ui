package http

import (
	"crypto/tls"
	"net"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/boreq/bolt-ui/internal/config"
	"github.com/boreq/bolt-ui/logging"
	"github.com/boreq/errors"
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

func (s *Server) Serve() error {
	handler := s.handler

	if s.conf.InsecureCORS {
		handler = cors.AllowAll().Handler(s.handler)
	}

	handler = gziphandler.GzipHandler(handler)

	if s.conf.InsecureTLS {
		s.log.Debug("starting an insecure listener", "address", s.conf.ServeAddress)
		return http.ListenAndServe(s.conf.ServeAddress, handler)
	}

	s.log.Debug("starting listening", "address", s.conf.ServeAddress)

	l, err := net.Listen("tcp", s.conf.ServeAddress)
	if err != nil {
		return errors.Wrap(err, "could not create listener")
	}

	l = tls.NewListener(l, &tls.Config{
		Certificates: []tls.Certificate{
			s.conf.Certificate,
		},
	})

	return http.Serve(l, handler)

}
