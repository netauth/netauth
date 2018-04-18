package rpc

import (
	"github.com/NetAuth/NetAuth/internal/token"
)

type EntityTree interface {
	ValidateSecret(string, string) error
	MakeBootstrap(string, string)
	DisableBootstrap()
}

type NetAuthServer struct {
	Tree EntityTree
	Token token.TokenService
}
