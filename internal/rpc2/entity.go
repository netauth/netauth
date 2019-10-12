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

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "EntityCreate",
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	// Token validation and authorization
	var err error
	ctx, err = s.checkToken(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.isAuthorized(ctx, types.Capability_CREATE_ENTITY); err != nil {
		return &pb.Empty{}, err
	}

	switch err := s.CreateEntity(e.GetID(), e.GetNumber(), e.GetSecret()); err {
	case tree.ErrDuplicateEntityID, tree.ErrDuplicateNumber:
		s.log.Warn("Attempt to create duplicate entity",
			"entity", e.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrExists
	case nil:
		s.log.Info("Entity Created",
			"entity", e.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, nil
	default:
		s.log.Warn("Error Creating Entity",
			"entity", e.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
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
	de := r.GetData()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "EntityUpdate",
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	// Token validation and authorization
	var err error
	ctx, err = s.checkToken(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.isAuthorized(ctx, types.Capability_MODIFY_ENTITY_META); err != nil {
		return &pb.Empty{}, err
	}

	switch err := s.UpdateEntityMeta(de.GetID(), de.GetMeta()); err {
	case db.ErrUnknownEntity:
		s.log.Warn("Entity does not exist!",
			"method", "EntityUpdate",
			"entity", de.GetID(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.Empty{}, ErrDoesNotExist

	default:
		s.log.Warn("Error Updating Entity",
			"entity", de.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	case nil:
		s.log.Info("Entity Updated",
			"entity", de.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, nil
	}
}

// EntityInfo provides information on a single entity.  The list
// returned is guaranteed to be of length 1.
func (s *Server) EntityInfo(ctx context.Context, r *pb.EntityRequest) (*pb.ListOfEntities, error) {
	e := r.GetEntity()

	switch ent, err := s.FetchEntity(e.GetID()); err {
	case db.ErrUnknownEntity:
		s.log.Warn("Entity does not exist!",
			"method", "EntityUpdate",
			"entity", e.GetID(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfEntities{}, ErrDoesNotExist
	default:
		s.log.Warn("Error fetching entity",
			"entity", e.GetID(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfEntities{}, ErrInternal
	case nil:
		s.log.Info("Dumped Entity Info",
			"entity", e.GetID(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfEntities{Entities: []*types.Entity{ent}}, nil
	}
}

// EntitySearch searches all entities and returns the entities that
// had been found.
func (s *Server) EntitySearch(ctx context.Context, r *pb.SearchRequest) (*pb.ListOfEntities, error) {
	expr := r.GetExpression()

	res, err := s.SearchEntities(db.SearchRequest{Expression: expr})
	if err != nil {
		s.log.Warn("Search Error",
			"expr", expr,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfEntities{}, ErrInternal
	}

	return &pb.ListOfEntities{Entities: res}, nil
}

// EntityUM handles both updates, and reads to the untyped metadata
// that's stored on Entities.
func (s *Server) EntityUM(ctx context.Context, r *pb.KVRequest) (*pb.ListOfStrings, error) {
	if r.GetAction() != pb.Action_READ &&
		r.GetAction() != pb.Action_UPSERT &&
		r.GetAction() != pb.Action_CLEAREXACT &&
		r.GetAction() != pb.Action_CLEARFUZZY {
		return &pb.ListOfStrings{}, ErrMalformedRequest
	}

	if r.GetAction() != pb.Action_READ {
		if s.readonly {
			s.log.Warn("Mutable request in read-only mode!",
				"method", "EntityUM",
				"client", getClientName(ctx),
				"service", getServiceName(ctx),
			)
			return &pb.ListOfStrings{}, ErrReadOnly
		}

		// Token validation and authorization
		var err error
		ctx, err = s.checkToken(ctx)
		if err != nil {
			return &pb.ListOfStrings{}, err
		}
		if err := s.isAuthorized(ctx, types.Capability_MODIFY_ENTITY_META); err != nil {
			return &pb.ListOfStrings{}, err
		}
	}

	// At this point, we're either in a read-only query, or in a
	// write one that has been authorized.
	meta, err := s.ManageUntypedEntityMeta(r.GetTarget(), r.GetAction().String(), r.GetKey(), r.GetValue())
	switch err {
	case db.ErrUnknownEntity:
		s.log.Warn("Entity does not exist!",
			"method", "EntityUM",
			"entity", r.GetTarget(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfStrings{}, ErrDoesNotExist

	default:
		s.log.Warn("Error Updating Entity",
			"entity", r.GetTarget(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfStrings{}, ErrInternal
	case nil:
		s.log.Info("Entity Updated",
			"entity", r.GetTarget(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfStrings{Strings: meta}, nil
	}
}

// EntityKeys handles updates and reads to keys for entities.
func (s *Server) EntityKeys(ctx context.Context, r *pb.KVRequest) (*pb.ListOfStrings, error) {
	if r.GetAction() != pb.Action_READ &&
		r.GetAction() != pb.Action_ADD &&
		r.GetAction() != pb.Action_DROP {
		return &pb.ListOfStrings{}, ErrMalformedRequest
	}

	if r.GetAction() != pb.Action_READ {
		if s.readonly {
			s.log.Warn("Mutable request in read-only mode!",
				"method", "EntityUM",
				"client", getClientName(ctx),
				"service", getServiceName(ctx),
			)
			return &pb.ListOfStrings{}, ErrReadOnly
		}

		// Token validation and authorization
		var err error
		ctx, err = s.checkToken(ctx)
		if err != nil {
			return &pb.ListOfStrings{}, err
		}
		err = s.isAuthorized(ctx, types.Capability_MODIFY_ENTITY_KEYS)
		if err != nil && getTokenClaims(ctx).EntityID != r.GetTarget() {
			return &pb.ListOfStrings{}, err
		}
	}

	// At this point, we're either in a read-only query, or in a
	// write one that has been authorized.
	keys, err := s.UpdateEntityKeys(r.GetTarget(), r.GetAction().String(), r.GetKey(), r.GetValue())
	switch err {
	case db.ErrUnknownEntity:
		s.log.Warn("Entity does not exist!",
			"method", "EntityUM",
			"entity", r.GetTarget(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfStrings{}, ErrDoesNotExist

	default:
		s.log.Warn("Error Updating Entity",
			"entity", r.GetTarget(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfStrings{}, ErrInternal
	case nil:
		s.log.Info("Entity Updated",
			"entity", r.GetTarget(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfStrings{Strings: keys}, nil
	}
}

// EntityDestroy will remove an entity from the system.  This is
// generally discouraged, but if you must then this function will do
// it.
func (s *Server) EntityDestroy(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	e := r.GetEntity()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "EntityDestroy",
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	// Token validation and authorization
	var err error
	ctx, err = s.checkToken(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.isAuthorized(ctx, types.Capability_DESTROY_ENTITY); err != nil {
		return &pb.Empty{}, err
	}

	switch err := s.DestroyEntity(e.GetID()); err {
	case db.ErrUnknownEntity:
		s.log.Warn("Entity does not exist!",
			"method", "EntityDestroy",
			"entity", e.GetID(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.Empty{}, ErrDoesNotExist

	default:
		s.log.Warn("Error Updating Entity",
			"entity", e.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	case nil:
		s.log.Info("Entity Updated",
			"entity", e.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, nil
	}
}

// EntityLock sets the lock flag on an entity.
func (s *Server) EntityLock(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	e := r.GetEntity()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "EntityLock",
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	// Token validation and authorization
	var err error
	ctx, err = s.checkToken(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.isAuthorized(ctx, types.Capability_LOCK_ENTITY); err != nil {
		return &pb.Empty{}, err
	}

	switch err := s.LockEntity(e.GetID()); err {
	case db.ErrUnknownEntity:
		s.log.Warn("Entity does not exist!",
			"method", "EntityLock",
			"entity", e.GetID(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.Empty{}, ErrDoesNotExist

	default:
		s.log.Warn("Error Locking Entity",
			"entity", e.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	case nil:
		s.log.Info("Entity Locked",
			"entity", e.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, nil
	}
}

// EntityUnlock clears the lock flag on an entity.
func (s *Server) EntityUnlock(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	e := r.GetEntity()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "EntityUnlock",
			"client", getClientName(ctx),
			"service", getServiceName(ctx),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	// Token validation and authorization
	var err error
	ctx, err = s.checkToken(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.isAuthorized(ctx, types.Capability_UNLOCK_ENTITY); err != nil {
		return &pb.Empty{}, err
	}

	switch err := s.UnlockEntity(e.GetID()); err {
	case db.ErrUnknownEntity:
		s.log.Warn("Entity does not exist!",
			"method", "EntityUnlock",
			"entity", e.GetID(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.Empty{}, ErrDoesNotExist

	default:
		s.log.Warn("Error Unlocking Entity",
			"entity", e.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	case nil:
		s.log.Info("Entity Unlocked",
			"entity", e.GetID(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, nil
	}
}

// EntityGroups returns the full membership for a given entity.
func (s *Server) EntityGroups(ctx context.Context, r *pb.EntityRequest) (*pb.ListOfGroups, error) {
	e := r.GetEntity()

	ent, err := s.FetchEntity(e.GetID())
	switch err {
	case db.ErrUnknownEntity:
		s.log.Warn("Entity does not exist!",
			"method", "EntityGroups",
			"entity", e.GetID(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfGroups{}, ErrDoesNotExist

	default:
		s.log.Warn("Error getting groups for entity",
			"entity", e.GetID(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfGroups{}, ErrInternal
	case nil:
		break
	}

	// Summoning memberhips without indirects is confusing, not
	// fully understood by the general user, and of exceedingly
	// limited utility.  As a result of this the feature is being
	// removed.  Until that happens though, this interface still
	// needs to be fed a "true" to maintain the feature.
	groups := s.GetMemberships(ent, true)

	out := make([]*types.Group, len(groups))
	for i := range groups {
		// We throw this error out here, as its logged at a
		// lower level, and the side effect here is that only
		// a partial result gets returned.
		tmp, _ := s.FetchGroup(groups[i])
		out[i] = tmp
	}

	return &pb.ListOfGroups{Groups: out}, nil
}
