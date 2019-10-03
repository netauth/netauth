package rpc2

import (
	"context"

	"github.com/NetAuth/NetAuth/internal/tree"
	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol/v2"
	types "github.com/NetAuth/Protocol"
)

// GroupCreate provisions a new group on the system.
func (s *Server) GroupCreate(ctx context.Context, r *pb.GroupRequest) (*pb.Empty, error) {
	g := r.GetGroup()
	client := r.GetInfo()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "GroupCreate",
			"client", client.GetID(),
			"service", client.GetService(),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	c, err := s.Validate(r.GetAuth().GetToken())
	if err != nil {
		s.log.Info("Permission Denied",
			"method", "GroupCreate",
			"group", g.GetName(),
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrUnauthenticated
	}
	if !c.HasCapability(types.Capability_CREATE_GROUP) {
		s.log.Info("Permission Denied",
			"method", "GroupCreate",
			"group", g.GetName(),
			"service", client.GetService(),
			"client", client.GetID(),
			"error", "missing-capability",
		)
		return &pb.Empty{}, ErrRequestorUnqualified
	}

	switch err := s.CreateGroup(g.GetName(), g.GetDisplayName(), g.GetManagedBy(), g.GetNumber()); err {
	case tree.ErrDuplicateGroupName, tree.ErrDuplicateNumber:
		s.log.Warn("Attempt to create duplicate group",
			"group", g.GetName(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrExists
	case nil:
		s.log.Info("Group Created",
			"group", g.GetName(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, nil
	default:
		s.log.Warn("Error Creating Group",
			"group", g.GetName(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}
}

// GroupUpdate adjusts the metadata on a group with the exception of
// untyped metadata.
func (s *Server) GroupUpdate(ctx context.Context, r *pb.GroupRequest) (*pb.Empty, error) {
	g := r.GetGroup()
	client := r.GetInfo()

	if s.readonly {
		s.log.Warn("Mutable request in read-only mode!",
			"method", "GroupUpdate",
			"client", client.GetID(),
			"service", client.GetService(),
		)
		return &pb.Empty{}, ErrReadOnly
	}

	c, err := s.Validate(r.GetAuth().GetToken())
	if err != nil {
		s.log.Info("Permission Denied",
			"method", "GroupCreate",
			"group", g.GetName(),
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrUnauthenticated
	}
	if !c.HasCapability(types.Capability_MODIFY_GROUP_META) {
		s.log.Info("Permission Denied",
			"method", "GroupCreate",
			"group", g.GetName(),
			"service", client.GetService(),
			"client", client.GetID(),
			"error", "missing-capability",
		)
		return &pb.Empty{}, ErrRequestorUnqualified
	}

	switch err := s.UpdateGroupMeta(g.GetName(), g); err {
	case db.ErrUnknownGroup:
		s.log.Warn("Unable to load group",
			"group", g.GetName(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrDoesNotExist
	case nil:
		s.log.Info("Group Updated",
			"group", g.GetName(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, nil
	default:
		s.log.Warn("Error Updating Group",
			"group", g.GetName(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
			"error", err,
		)
		return &pb.Empty{}, ErrInternal
	}
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
