package rpc

import pb "github.com/NetAuth/NetAuth/proto"

type EntityManager interface {
	NewEntity(string, string, string, int32, string) error
	DeleteEntity(string, string, string) error
	ChangeSecret(string, string, string, string) error
	ValidateSecret(string, string) error
	MakeBootstrap(string, string)
	DisableBootstrap()
	ListMembers(string) ([]*pb.Entity, error)
}

type NetAuthServer struct {
	EM EntityManager
}
