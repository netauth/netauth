package rpc2

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/pkg/token"

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
	CreateEntity(context.Context, string, int32, string) error
	FetchEntity(context.Context, string) (*pb.Entity, error)
	SearchEntities(context.Context, db.SearchRequest) ([]*pb.Entity, error)
	ValidateSecret(context.Context, string, string) error
	SetSecret(context.Context, string, string) error
	LockEntity(context.Context, string) error
	UnlockEntity(context.Context, string) error
	UpdateEntityMeta(context.Context, string, *pb.EntityMeta) error
	EntityKVGet(context.Context, string, []*pb.KVData) ([]*pb.KVData, error)
	EntityKVAdd(context.Context, string, []*pb.KVData) error
	EntityKVDel(context.Context, string, []*pb.KVData) error
	EntityKVReplace(context.Context, string, []*pb.KVData) error
	UpdateEntityKeys(context.Context, string, string, string, string) ([]string, error)
	ManageUntypedEntityMeta(context.Context, string, string, string, string) ([]string, error)
	DestroyEntity(context.Context, string) error

	CreateGroup(context.Context, string, string, string, int32) error
	FetchGroup(context.Context, string) (*pb.Group, error)
	SearchGroups(context.Context, db.SearchRequest) ([]*pb.Group, error)
	UpdateGroupMeta(context.Context, string, *pb.Group) error
	ManageUntypedGroupMeta(context.Context, string, string, string, string) ([]string, error)
	GroupKVGet(context.Context, string, []*pb.KVData) ([]*pb.KVData, error)
	GroupKVAdd(context.Context, string, []*pb.KVData) error
	GroupKVDel(context.Context, string, []*pb.KVData) error
	GroupKVReplace(context.Context, string, []*pb.KVData) error
	DestroyGroup(context.Context, string) error

	AddEntityToGroup(context.Context, string, string) error
	RemoveEntityFromGroup(context.Context, string, string) error
	ListMembers(context.Context, string) ([]*pb.Entity, error)
	GetMemberships(context.Context, *pb.Entity) []string
	ModifyGroupRule(context.Context, string, string, rpc.RuleAction) error

	SetEntityCapability2(context.Context, string, *pb.Capability) error
	DropEntityCapability2(context.Context, string, *pb.Capability) error
	SetGroupCapability2(context.Context, string, *pb.Capability) error
	DropGroupCapability2(context.Context, string, *pb.Capability) error
}

// Options configure the server
type Option func(s *Server)
