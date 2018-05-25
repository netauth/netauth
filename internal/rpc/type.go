package rpc

import (
	"github.com/NetAuth/NetAuth/internal/token"

	pb "github.com/NetAuth/Protocol"
)

type EntityTree interface {
	GetEntity(string) (*pb.Entity, error)
	ValidateSecret(string, string) error
	MakeBootstrap(string, string)
	DisableBootstrap()
	SetEntitySecretByID(string, string) error

	NewEntity(string, int32, string) error
	DeleteEntityByID(string) error
	UpdateEntityMeta(string, *pb.EntityMeta) error

	NewGroup(string, string, string, int32) error
	DeleteGroup(string) error
	ListGroups() ([]*pb.Group, error)
	UpdateGroupMeta(string, *pb.Group) error

	AddEntityToGroup(string, string) error
	RemoveEntityFromGroup(string, string) error
	ListMembers(string) ([]*pb.Entity, error)
}

type NetAuthServer struct {
	Tree  EntityTree
	Token token.TokenService
}
