package rpc2

import (
	"context"

	pb "github.com/NetAuth/Protocol/v2"
)

// EntityCreate creates entities.  This call will validate that a
// correct token is held, which must contain either CREATE_ENTITY or
// GLOBAL_ROOT permissions.
func (s *Server) EntityCreate(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// EntityUpdate provides a change to specific entity metadata that is
// in the typed data fields.  This method does not update keys,
// groups, untyped metadata, or capabilities.  To call this method you
// must be in posession of a token with MODIFY_ENTITY_META
// capabilities.
func (s *Server) EntityUpdate(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// EntityInfo provides information on a single entity.  The list
// returned is guaranteed to be of length 1.
func (s *Server) EntityInfo(ctx context.Context, r *pb.EntityRequest) (*pb.ListOfEntities, error) {
	return &pb.ListOfEntities{}, nil
}

// EntitySearch searches all entities and returns the entities that
// had been found.
func (s *Server) EntitySearch(ctx context.Context, r *pb.SearchRequest) (*pb.ListOfEntities, error) {
	return &pb.ListOfEntities{}, nil
}

// EntityUM handles both updates, and reads to the untyped metadata
// that's stored on Entities.
func (s *Server) EntityUM(ctx context.Context, r *pb.KVRequest) (*pb.ListOfStrings, error) {
	return &pb.ListOfStrings{}, nil
}

// EntityKeys handles updates and reads to keys for entities.
func (s *Server) EntityKeys(ctx context.Context, r *pb.KVRequest) (*pb.ListOfStrings, error) {
	return &pb.ListOfStrings{}, nil
}

// EntityDestroy will remove an entity from the system.  This is
// generally discouraged, but if you must then this function will do
// it.
func (s *Server) EntityDestroy(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// EntityLock sets the lock flag on an entity.
func (s *Server) EntityLock(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// EntityUnlock clears the lock flag on an entity.
func (s *Server) EntityUnlock(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

// EntityGroups returns the full membership for a given entity.
func (s *Server) EntityGroups(ctx context.Context, r *pb.EntityRequest) (*pb.ListOfGroups, error) {
	return &pb.ListOfGroups{}, nil
}
