package rpc

import (
	"context"
	"log"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

// SearchEntities allows searches to be run for entities.
func (s *NetAuthServer) SearchEntities(ctx context.Context, r *pb.SearchRequest) (*pb.EntityList, error) {
	srchexpr := r.GetExpression()
	client := r.GetInfo()

	res, err := s.Tree.SearchEntities(db.SearchRequest{Expression: srchexpr})
	if err != nil {
		return nil, toWireError(err)
	}

	log.Printf("Entity Search '%s' (%s@%s)", srchexpr, client.GetService(), client.GetID())
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

	log.Printf("Group Search '%s' (%s@%s)", srchexpr, client.GetService(), client.GetID())
	return &pb.GroupList{Groups: res}, toWireError(nil)
}
