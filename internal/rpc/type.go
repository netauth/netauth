package rpc

import (
	"github.com/hashicorp/go-hclog"

	"github.com/NetAuth/NetAuth/internal/token"
	"github.com/NetAuth/NetAuth/internal/rpc2"
)

// An EntityTree is a mechanism for storing entities and information
// about them.
type EntityTree = rpc2.Manager

// A NetAuthServer is a collection of methods that satisfy the
// requirements of the NetAuthServer protocol buffer.
type NetAuthServer struct {
	Tree  EntityTree
	Token token.Service
	Log   hclog.Logger
}
