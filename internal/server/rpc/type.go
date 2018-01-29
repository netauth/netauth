package rpc

import pb "github.com/NetAuth/NetAuth/pkg/proto"

type EntityManager interface {
	NewEntity(string, string, string, int32, string) error
	DeleteEntity(string, string, string) error
	ChangeSecret(string, string, string, string) error
	ValidateSecret(string, string) error
	MakeBootstrap(string, string)
	DisableBootstrap()
	ListMembers(string) ([]*pb.Entity, error)
	GetEntity(string) (*pb.Entity, error)
	UpdateEntityMeta(string, string, string, *pb.EntityMeta) error
	NewGroup(string, string, string, string, int32) error
	DeleteGroup(string, string, string) error
	UpdateGroupMeta(string, string, string, *pb.Group) error
}

type NetAuthServer struct {
	EM EntityManager
}
