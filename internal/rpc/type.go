package rpc

import (
	"github.com/NetAuth/NetAuth/internal/token"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

type EntityTree interface {
	GetEntity(string) (*pb.Entity, error)
	ValidateSecret(string, string) error
	MakeBootstrap(string, string)
	DisableBootstrap()
	SetEntitySecretByID(string, string) error
}

type NetAuthServer struct {
	Tree EntityTree
	Token token.TokenService
}
