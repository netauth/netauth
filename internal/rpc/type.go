package rpc

import (
	"errors"

	"github.com/NetAuth/NetAuth/internal/token"

	pb "github.com/NetAuth/Protocol"
)

var (
	RequestorUnqualified = errors.New("The requestor is not qualified to perform that action")
	MalformedRequest     = errors.New("The request is malformed and cannot be processed")
	InternalError        = errors.New("An internal error has occured")
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
	UpdateEntityKeys(string, string, string, string) ([]string, error)

	NewGroup(string, string, string, int32) error
	DeleteGroup(string) error
	ListGroups() ([]*pb.Group, error)
	GetGroupByName(string) (*pb.Group, error)
	UpdateGroupMeta(string, *pb.Group) error
	GetMemberships(*pb.Entity, bool) []string

	AddEntityToGroup(string, string) error
	RemoveEntityFromGroup(string, string) error
	ListMembers(string) ([]*pb.Entity, error)

	ModifyGroupExpansions(string, string, pb.ExpansionMode) error

	SetEntityCapabilityByID(string, string) error
	RemoveEntityCapabilityByID(string, string) error
	SetGroupCapabilityByName(string, string) error
	RemoveGroupCapabilityByName(string, string) error
}

type NetAuthServer struct {
	Tree  EntityTree
	Token token.TokenService
}
