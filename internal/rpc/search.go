package rpc

import (
	"context"

	"github.com/netauth/netauth/internal/db"

	pb "github.com/netauth/protocol"
)

// SearchEntities allows searches to be run for entities.
func (s *NetAuthServer) SearchEntities(ctx context.Context, r *pb.SearchRequest) (*pb.EntityList, error) {
	srchexpr := r.GetExpression()
	client := r.GetInfo()

	res, err := s.Tree.SearchEntities(db.SearchRequest{Expression: srchexpr})
	if err != nil {
		return nil, toWireError(err)
	}

	s.Log.Info("Entity Search", "expression", srchexpr, "service", client.GetService(), "client", client.GetID())
	return &pb.EntityList{Members: res}, toWireError(nil)
}

// SearchGroups allows searches to be run for groups.
func (s *NetAuthServer) SearchGroups(ctx context.Context, r *pb.SearchRequest) (*pb.GroupList, error) {
	srchexpr := r.GetExpression()
	client := r.GetInfo()

	res, err := s.Tree.SearchGroups(db.SearchRequest{Expression: srchexpr})
	if err != nil {
		return nil, toWireError(err)
	}

	s.Log.Info("Group Search", "expression", srchexpr, "service", client.GetService(), "client", client.GetID())
	return &pb.GroupList{Groups: res}, toWireError(nil)
}
