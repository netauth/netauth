package rpc2

import (
	"context"

	"github.com/NetAuth/NetAuth/internal/db"
	"github.com/NetAuth/NetAuth/internal/tree"

	types "github.com/NetAuth/Protocol"
	pb "github.com/NetAuth/Protocol/v2"
)

// EntityCreate creates entities.  This call will validate that a
// correct token is held, which must contain either CREATE_ENTITY or
// GLOBAL_ROOT permissions.
func (s *Server) EntityCreate(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	e := r.GetEntity()
	client := r.GetInfo()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "EntityCreate",
			"client", client.GetID(),
			"service", client.GetService(),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	c, err := s.Validate(r.GetAuth().GetToken())
	if err != nil {
		s.log.Info("Permission Denied",
			"method", "EntityCreate",
			"entity", e.GetID(),
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrUnauthenticated
	}
	if !c.HasCapability(types.Capability_CREATE_ENTITY) {
		s.log.Info("Permission Denied",
			"method", "EntityCreate",
			"entity", e.GetID(),
			"service", client.GetService(),
			"client", client.GetID(),
			"error", "missing-capability",
		)
		return &pb.Empty{}, ErrRequestorUnqualified
	}

	switch err := s.CreateEntity(e.GetID(), e.GetNumber(), e.GetSecret()); err {
	case tree.ErrDuplicateEntityID, tree.ErrDuplicateNumber:
		s.log.Warn("Attempt to create duplicate entity",
			"entity", e.GetID(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrExists
	case nil:
		s.log.Info("Entity Created",
			"entity", e.GetID(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, nil
	default:
		s.log.Warn("Error Creating Entity",
			"entity", e.GetID(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}
}

// EntityUpdate provides a change to specific entity metadata that is
// in the typed data fields.  This method does not update keys,
// groups, untyped metadata, or capabilities.  To call this method you
// must be in posession of a token with MODIFY_ENTITY_META
// capabilities.
func (s *Server) EntityUpdate(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	client := r.GetInfo()
	de := r.GetData()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "EntityUpdate",
			"client", client.GetID(),
			"service", client.GetService(),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	c, err := s.Validate(r.GetAuth().GetToken())
	if err != nil {
		s.log.Info("Permission Denied",
			"method", "EntityUpdate",
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrUnauthenticated
	}
	if !c.HasCapability(types.Capability_MODIFY_ENTITY_META) {
		s.log.Info("Permission Denied",
			"method", "EntityUpdate",
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", "missing-capability",
		)
		return &pb.Empty{}, ErrRequestorUnqualified
	}

	switch err := s.UpdateEntityMeta(de.GetID(), de.GetMeta()); err {
	case db.ErrUnknownEntity:
		s.log.Warn("Entity does not exist!",
			"method", "EntityUpdate",
			"entity", de.GetID(),
			"service", client.GetService(),
			"client", client.GetID(),
		)
		return &pb.Empty{}, ErrDoesNotExist

	default:
		s.log.Warn("Error Updating Entity",
			"entity", de.GetID(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	case nil:
		s.log.Info("Entity Updated",
			"entity", de.GetID(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, nil
	}
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
