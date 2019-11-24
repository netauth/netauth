package rpc2

import (
	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/token"

	pb "github.com/netauth/protocol"
	rpc "github.com/netauth/protocol/v2"
)

// Server returns the interface which satisfies the gRPC type for the
// server.
type Server struct {
	token.Service
	Manager

	readonly bool
	log      hclog.Logger
}

// Refs is the container that is used to provide references to the RPC
// server.
type Refs struct {
	TokenService token.Service
	Tree         Manager
}

// The Manager handles backend data and is an equivalent interface to rpc.EntityTree
type Manager interface {
	RegisterEntityHookToChain(string, string) error
	RegisterGroupHookToChain(string, string) error

	Bootstrap(string, string)
	DisableBootstrap()

	CreateEntity(string, int32, string) error
	FetchEntity(string) (*pb.Entity, error)
	SearchEntities(db.SearchRequest) ([]*pb.Entity, error)
	ValidateSecret(string, string) error
	SetSecret(string, string) error
	LockEntity(string) error
	UnlockEntity(string) error
	UpdateEntityMeta(string, *pb.EntityMeta) error
	UpdateEntityKeys(string, string, string, string) ([]string, error)
	ManageUntypedEntityMeta(string, string, string, string) ([]string, error)
	DestroyEntity(string) error

	CreateGroup(string, string, string, int32) error
	FetchGroup(string) (*pb.Group, error)
	SearchGroups(db.SearchRequest) ([]*pb.Group, error)
	UpdateGroupMeta(string, *pb.Group) error
	ManageUntypedGroupMeta(string, string, string, string) ([]string, error)
	DestroyGroup(string) error

	AddEntityToGroup(string, string) error
	RemoveEntityFromGroup(string, string) error
	ListMembers(string) ([]*pb.Entity, error)
	GetMemberships(*pb.Entity, bool) []string
	ModifyGroupRule(string, string, rpc.RuleAction) error

	SetEntityCapability2(string, *pb.Capability) error
	DropEntityCapability2(string, *pb.Capability) error
	SetGroupCapability2(string, *pb.Capability) error
	DropGroupCapability2(string, *pb.Capability) error
}
