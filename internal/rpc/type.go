package rpc

import (
	"errors"

	"github.com/NetAuth/NetAuth/internal/token"

	pb "github.com/NetAuth/Protocol"
)

var (
	// ErrRequestorUnqualified is returned when a caller has
	// attempted to perform some action that requires
	// authorization and the caller is either not authorized, was
	// unable to present a token, or the token did not contain
	// sufficient capabilities.
	ErrRequestorUnqualified = errors.New("the requestor is not qualified to perform that action")

	// ErrMalformedRequest is returned when a caller makes some
	// request to the server but has failed to provide a complete
	// request, or has provided a request that is in conflict with
	// itself.
	ErrMalformedRequest = errors.New("the request is malformed and cannot be processed")

	// ErrInternalError is a catchall for errors that are
	// otherwise unidentified and unrecoverable in the server.
	ErrInternalError = errors.New("An internal error has occured")
)

// An EntityTree is a mechanism for storing entities and information
// about them.
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

// A NetAuthServer is a collection of methods that satisfy the
// requirements of the NetAuthServer protocol buffer.
type NetAuthServer struct {
	Tree  EntityTree
	Token token.Service
}
