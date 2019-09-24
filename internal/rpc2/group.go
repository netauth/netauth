package rpc2

import (
	"context"

	pb "github.com/NetAuth/Protocol/v2"
)

// GroupCreate provisions a new group on the system.
func (s *Server) GroupCreate(ctx context.Context, r *pb.GroupRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// GroupUpdate adjusts the metadata on a group with the exception of
// untyped metadata.
func (s *Server) GroupUpdate(ctx context.Context, r *pb.GroupRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// GroupInfo returns a group for inspection.  It does not return
// key/value data.
func (s *Server) GroupInfo(ctx context.Context, r *pb.GroupRequest) (*pb.ListOfGroups, error) {
	return &pb.ListOfGroups{}, nil
}

// GroupUM handles updates to untyped metadata for groups.
func (s *Server) GroupUM(ctx context.Context, r *pb.KVRequest) (*pb.ListOfStrings, error) {
	return &pb.ListOfStrings{}, nil
}

// GroupUpdateRules updates the expansion rules on a particular group.
func (s *Server) GroupUpdateRules(ctx context.Context, r *pb.GroupRulesRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// GroupAddMember adds an entity directly to a group.
func (s *Server) GroupAddMember(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// GroupDelMember dels an entity directly to a group.
func (s *Server) GroupDelMember(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// GroupDestroy will remove a group from the server completely.  This
// is not recommended and should not be done, but if you must here it
// is.
func (s *Server) GroupDestroy(ctx context.Context, r *pb.GroupRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// GroupMembers returns the list of all entities that are members of
// the group.
func (s *Server) GroupMembers(ctx context.Context, r *pb.GroupRequest) (*pb.ListOfEntities, error) {
	return &pb.ListOfEntities{}, nil
}

// GroupSearch searches for groups and returns a list of all groups
// matching the criteria specified.
func (s *Server) GroupSearch(ctx context.Context, r *pb.SearchRequest) (*pb.ListOfGroups, error) {
	return &pb.ListOfGroups{}, nil
}
