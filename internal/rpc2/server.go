package rpc2

import (
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/pkg/token"
)

// New returns a ready to use server implementation.
func New(opts ...Option) *Server {
	s := &Server{
		log:      hclog.NewNullLogger(),
		readonly: false,
	}

	for _, o := range opts {
		o(s)
	}
	return s
}

func WithLogger(l hclog.Logger) Option { return func(s *Server) { s.log = l } }

func WithTokenService(t token.Service) Option { return func(s *Server) { s.Service = t } }

func WithEntityTree(t Manager) Option { return func(s *Server) { s.Manager = t } }

func WithDisabledWrites(r bool) Option { return func(s *Server) { s.readonly = r } }
