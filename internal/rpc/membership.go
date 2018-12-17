package rpc

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

// AddEntityToGroup will add an existing entity to an existing group
// if they are not already a direct member.  If they are a direct
// member this call is idempotent.  This action must be authorized by
// the presentation of a token containing the appropriate capability.
func (s *NetAuthServer) AddEntityToGroup(ctx context.Context, r *pb.ModEntityMembershipRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	g := r.GetGroup()
	e := r.GetEntity()

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Either the entity must posses the right capability, or they
	// must be in the a group that is permitted to manage this one
	// based on membership.  Either is sufficient.
	if !s.manageByMembership(c.EntityID, g.GetName()) && !c.HasCapability("MODIFY_GROUP_MEMBERS") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	// Add to the group
	if err := s.Tree.AddEntityToGroup(e.GetID(), g.GetName()); err != nil {
		return nil, toWireError(err)
	}

	log.Printf("Entity '%s' added to '%s' by '%s' (%s@%s)",
		e.GetID(),
		g.GetName(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Membership updated successfully"),
		Success: proto.Bool(true),
	}, toWireError(nil)
}

// RemoveEntityFromGroup will remove an existing entity from an
// existing group.  This action must be authorized by the presentation
// of a token containing appropriate capabilities.
func (s *NetAuthServer) RemoveEntityFromGroup(ctx context.Context, r *pb.ModEntityMembershipRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	g := r.GetGroup()
	e := r.GetEntity()

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Either the entity must posses the right capability, or they
	// must be in the a group that is permitted to manage this one
	// based on membership.  Either is sufficient.
	if !s.manageByMembership(c.EntityID, g.GetName()) && !c.HasCapability("MODIFY_GROUP_MEMBERS") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	// Remove from the group
	if err := s.Tree.RemoveEntityFromGroup(e.GetID(), g.GetName()); err != nil {
		return nil, toWireError(err)
	}

	log.Printf("Entity '%s' removed from '%s' by '%s' (%s@%s)",
		e.GetID(),
		g.GetName(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Membership updated successfully"),
		Success: proto.Bool(true),
	}, toWireError(nil)
}

// ListGroups lists the groups a particular entity is in, or all
// groups on the server if no entity is specified.  In the case of
// calculating the groups a specific entity is in this can be quite
// expensive since large chunks of the membership tree will need to be
// calculated.
func (s *NetAuthServer) ListGroups(ctx context.Context, r *pb.GroupListRequest) (*pb.GroupList, error) {
	e := r.GetEntity()
	inclindr := r.GetIncludeIndirects()

	var list []*pb.Group

	if e != nil {
		// If e is defined then we want the groups for a
		// specific entity
		entity, err := s.Tree.FetchEntity(e.GetID())
		if err != nil {
			return nil, toWireError(err)
		}
		groupNames := s.Tree.GetMemberships(entity, inclindr)
		for _, name := range groupNames {
			g, err := s.Tree.FetchGroup(name)
			if err != nil {
				return nil, toWireError(err)
			}
			list = append(list, g)
		}
	} else {
		// If e is not defined then we want all groups.
		var err error
		list, err = s.Tree.SearchGroups(db.SearchRequest{Expression: "*"})
		if err != nil {
			return nil, toWireError(err)
		}
	}

	return &pb.GroupList{
		Groups: list,
	}, toWireError(nil)
}

// ListGroupMembers lists the members that are in a particular group.
// This call requires computing fairly large chunks of the membership
// graph.
func (s *NetAuthServer) ListGroupMembers(ctx context.Context, r *pb.GroupMemberRequest) (*pb.EntityList, error) {
	g := r.GetGroup()

	memberlist, err := s.Tree.ListMembers(g.GetName())
	if err != nil {
		return nil, toWireError(err)
	}

	return &pb.EntityList{
		Members: memberlist,
	}, toWireError(nil)
}

// ModifyGroupNesting permits changing the rules for group expansions.
// These expansions can either include a group's members, or prune the
// members of one group from another.  Expansions are checked to
// ensure they do not exist already, and that the addition of an
// expansion would not create a cycle in the membership graph.
func (s *NetAuthServer) ModifyGroupNesting(ctx context.Context, r *pb.ModGroupNestingRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	parent := r.GetParentGroup()
	child := r.GetChildGroup()
	mode := r.GetMode()

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Either the entity must posses the right capability, or they
	// must be in the a group that is permitted to manage this one
	// based on membership.  Either is sufficient.
	if !s.manageByMembership(c.EntityID, child.GetName()) && !c.HasCapability("MODIFY_GROUP_MEMBERS") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.ModifyGroupExpansions(parent.GetName(), child.GetName(), mode); err != nil {
		return nil, toWireError(err)
	}

	log.Printf("Group '%s'->'%s' expansion to '%s' by '%s' (%s@%s)",
		parent.GetName(),
		child.GetName(),
		mode,
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Nesting updated successfully"),
		Success: proto.Bool(true),
	}, toWireError(nil)
}
