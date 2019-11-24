package rpc2

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/netauth/netauth/internal/db"
	"github.com/netauth/netauth/internal/tree"

	types "github.com/netauth/protocol"
	pb "github.com/netauth/protocol/v2"
)

// GroupCreate provisions a new group on the system.
func (s *Server) GroupCreate(ctx context.Context, r *pb.GroupRequest) (*pb.Empty, error) {
	g := r.GetGroup()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "GroupCreate",
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
	if err := s.isAuthorized(ctx, types.Capability_CREATE_GROUP); err != nil {
		return &pb.Empty{}, err
	}

	switch err := s.CreateGroup(g.GetName(), g.GetDisplayName(), g.GetManagedBy(), g.GetNumber()); err {
	case tree.ErrDuplicateGroupName, tree.ErrDuplicateNumber:
		s.log.Warn("Attempt to create duplicate group",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrExists
	case nil:
		s.log.Info("Group Created",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, nil
	default:
		s.log.Warn("Error Creating Group",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}
}

// GroupUpdate adjusts the metadata on a group with the exception of
// untyped metadata.
func (s *Server) GroupUpdate(ctx context.Context, r *pb.GroupRequest) (*pb.Empty, error) {
	g := r.GetGroup()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "GroupUpdate",
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
	if err := s.isAuthorized(ctx, types.Capability_MODIFY_GROUP_META); err != nil && !s.manageByMembership(getTokenClaims(ctx).EntityID, g) {
		return &pb.Empty{}, err
	}

	switch err := s.UpdateGroupMeta(g.GetName(), g); err {
	case db.ErrUnknownGroup:
		s.log.Warn("Unable to load group",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrDoesNotExist
	case nil:
		s.log.Info("Group Updated",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, nil
	default:
		s.log.Warn("Error Updating Group",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}
}

// GroupInfo returns a group for inspection.  It does not return
// key/value data.
func (s *Server) GroupInfo(ctx context.Context, r *pb.GroupRequest) (*pb.ListOfGroups, error) {
	g := r.GetGroup()

	switch grp, err := s.FetchGroup(g.GetName()); err {
	case db.ErrUnknownGroup:
		s.log.Warn("Unknown Group",
			"group", g.GetName(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfGroups{}, ErrDoesNotExist
	case nil:
		s.log.Info("Group Info",
			"group", g.GetName(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfGroups{Groups: []*types.Group{grp}}, nil
	default:
		s.log.Warn("Error Loading Group",
			"group", g.GetName(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfGroups{}, ErrInternal
	}
}

// GroupUM handles updates to untyped metadata for groups.
func (s *Server) GroupUM(ctx context.Context, r *pb.KVRequest) (*pb.ListOfStrings, error) {
	if r.GetAction() != pb.Action_READ &&
		r.GetAction() != pb.Action_UPSERT &&
		r.GetAction() != pb.Action_CLEAREXACT &&
		r.GetAction() != pb.Action_CLEARFUZZY {
		return &pb.ListOfStrings{}, ErrMalformedRequest
	}

	if r.GetAction() != pb.Action_READ {
		if s.readonly {
			s.log.Warn("Mutable request in read-only mode!",
				"method", "GroupUM",
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
		g := types.Group{Name: proto.String(r.GetTarget())}
		if err := s.isAuthorized(ctx, types.Capability_MODIFY_GROUP_META); err != nil && !s.manageByMembership(getTokenClaims(ctx).EntityID, &g) {
			return &pb.ListOfStrings{}, err
		}
	}

	// At this point, we're either in a read-only query, or in a
	// write one that has been authorized.
	meta, err := s.ManageUntypedGroupMeta(r.GetTarget(), r.GetAction().String(), r.GetKey(), r.GetValue())
	switch err {
	case db.ErrUnknownGroup:
		s.log.Warn("Group does not exist!",
			"method", "GroupUM",
			"group", r.GetTarget(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfStrings{}, ErrDoesNotExist
	case nil:
		s.log.Info("Group Updated",
			"group", r.GetTarget(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfStrings{Strings: meta}, nil
	default:
		s.log.Warn("Error Updating Group",
			"group", r.GetTarget(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfStrings{}, ErrInternal
	}
}

// GroupUpdateRules updates the expansion rules on a particular group.
func (s *Server) GroupUpdateRules(ctx context.Context, r *pb.GroupRulesRequest) (*pb.Empty, error) {
	g := r.GetGroup()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "GroupUM",
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
	if err := s.isAuthorized(ctx, types.Capability_MODIFY_GROUP_MEMBERS); err != nil && !s.manageByMembership(getTokenClaims(ctx).EntityID, g) {
		return &pb.Empty{}, err
	}

	switch err := s.ModifyGroupRule(r.GetGroup().GetName(), r.GetTarget().GetName(), r.GetRuleAction()); err {
	case db.ErrUnknownGroup:
		s.log.Warn("Group does not exist!",
			"method", "GroupUpdateRules",
			"group", g.GetName(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.Empty{}, ErrDoesNotExist
	case nil:
		s.log.Info("Group Updated",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, nil
	default:
		s.log.Warn("Error Updating Group",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}
}

// GroupAddMember adds an entity directly to a group.
func (s *Server) GroupAddMember(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	e := r.GetEntity()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "GroupAddMember",
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

	for _, g := range e.GetMeta().GetGroups() {
		grp := types.Group{Name: proto.String(g)}
		if err := s.isAuthorized(ctx, types.Capability_MODIFY_GROUP_MEMBERS); err != nil && !s.manageByMembership(getTokenClaims(ctx).EntityID, &grp) {
			s.log.Warn("Insufficient authority to add entity to group",
				"entity", e.GetID(),
				"group", g,
				"authority", getTokenClaims(ctx).EntityID,
			)
			return &pb.Empty{}, err
		}
		if err := s.AddEntityToGroup(e.GetID(), g); err != nil {
			s.log.Warn("Error adding entity to group",
				"entity", e.GetID(),
				"group", g,
				"authority", getTokenClaims(ctx).EntityID,
				"service", getServiceName(ctx),
				"client", getClientName(ctx),
				"error", err,
			)
			return &pb.Empty{}, ErrInternal
		}
	}
	return &pb.Empty{}, nil
}

// GroupDelMember dels an entity directly to a group.
func (s *Server) GroupDelMember(ctx context.Context, r *pb.EntityRequest) (*pb.Empty, error) {
	e := r.GetEntity()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "GroupAddMember",
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

	for _, g := range e.GetMeta().GetGroups() {
		grp := types.Group{Name: proto.String(g)}
		if err := s.isAuthorized(ctx, types.Capability_MODIFY_GROUP_MEMBERS); err != nil && !s.manageByMembership(getTokenClaims(ctx).EntityID, &grp) {
			s.log.Warn("Insufficient authority to add entity to group",
				"entity", e.GetID(),
				"group", g,
				"authority", getTokenClaims(ctx).EntityID,
			)
			return &pb.Empty{}, err
		}
		if err := s.RemoveEntityFromGroup(e.GetID(), g); err != nil {
			s.log.Warn("Error adding entity to group",
				"entity", e.GetID(),
				"group", g,
				"authority", getTokenClaims(ctx).EntityID,
				"service", getServiceName(ctx),
				"client", getClientName(ctx),
				"error", err,
			)
			return &pb.Empty{}, ErrInternal
		}
	}
	return &pb.Empty{}, nil
}

// GroupDestroy will remove a group from the server completely.  This
// is not recommended and should not be done, but if you must here it
// is.
func (s *Server) GroupDestroy(ctx context.Context, r *pb.GroupRequest) (*pb.Empty, error) {
	g := r.GetGroup()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "GroupDestroy",
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
	if err := s.isAuthorized(ctx, types.Capability_DESTROY_GROUP); err != nil {
		return &pb.Empty{}, err
	}

	switch err := s.DestroyGroup(g.GetName()); err {
	case db.ErrUnknownGroup:
		s.log.Warn("Group does not exist!",
			"method", "GroupDestroy",
			"group", g.GetName(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.Empty{}, ErrDoesNotExist
	case nil:
		s.log.Info("Group Updated",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, nil
	default:
		s.log.Warn("Error Updating Group",
			"group", g.GetName(),
			"authority", getTokenClaims(ctx).EntityID,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}
}

// GroupMembers returns the list of all entities that are members of
// the group.
func (s *Server) GroupMembers(ctx context.Context, r *pb.GroupRequest) (*pb.ListOfEntities, error) {
	g := r.GetGroup()

	members, err := s.ListMembers(g.GetName())
	switch err {
	case db.ErrUnknownGroup:
		s.log.Warn("Group does not exist!",
			"method", "GroupDestroy",
			"group", g.GetName(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
		)
		return &pb.ListOfEntities{}, ErrDoesNotExist
	case nil:
		return &pb.ListOfEntities{Entities: members}, nil
	default:
		s.log.Warn("Error Fetching Membership Group",
			"group", g.GetName(),
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfEntities{}, ErrInternal
	}
}

// GroupSearch searches for groups and returns a list of all groups
// matching the criteria specified.
func (s *Server) GroupSearch(ctx context.Context, r *pb.SearchRequest) (*pb.ListOfGroups, error) {
	expr := r.GetExpression()

	res, err := s.SearchGroups(db.SearchRequest{Expression: expr})
	if err != nil {
		s.log.Warn("Search Error",
			"expr", expr,
			"service", getServiceName(ctx),
			"client", getClientName(ctx),
			"error", err,
		)
		return &pb.ListOfGroups{}, ErrInternal

	}
	return &pb.ListOfGroups{Groups: res}, nil
}
