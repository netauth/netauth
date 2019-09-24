package rpc2

import (
	"github.com/NetAuth/NetAuth/internal/token"
	"github.com/NetAuth/NetAuth/internal/tree"
)

// Server returns the interface which satisfies the gRPC type for the
// server.
type Server struct {
	token.Service
	tree.Manager
}

// Refs is the container that is used to provide references to the RPC
// server.
type Refs struct {
	TokenService token.Service
	Tree         tree.Manager
}

// New returns a ready to use server implementation.
func New(r Refs) *Server {
	return &Server{
		Service: r.TokenService,
		Manager: r.Tree,
	}
}
