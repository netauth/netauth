package rpc

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

func (s *NetAuthServer) AddEntityToGroup(ctx context.Context, r *pb.ModEntityMembershipRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	g := r.GetGroup()
	e := r.GetEntity()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Either the entity must posses the right capability, or they
	// must be in the a group that is permitted to manage this one
	// based on membership.  Either is sufficient.
	if !s.manageByMembership(c.EntityID, g.GetName()) && !c.HasCapability("MODIFY_GROUP_MEMBERS") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	// Add to the group
	if err := s.Tree.AddEntityToGroup(e.GetID(), g.GetName()); err != nil {
		return nil, err
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
	}, nil
}

func (s *NetAuthServer) RemoveEntityFromGroup(ctx context.Context, r *pb.ModEntityMembershipRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	g := r.GetGroup()
	e := r.GetEntity()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Either the entity must posses the right capability, or they
	// must be in the a group that is permitted to manage this one
	// based on membership.  Either is sufficient.
	if !s.manageByMembership(c.EntityID, g.GetName()) && !c.HasCapability("MODIFY_GROUP_MEMBERS") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	// Remove from the group
	if err := s.Tree.RemoveEntityFromGroup(e.GetID(), g.GetName()); err != nil {
		return nil, err
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
	}, nil
}

func (s *NetAuthServer) ListGroupMembers(ctx context.Context, r *pb.GroupMemberRequest) (*pb.EntityList, error) {
	client := r.GetInfo()
	g := r.GetGroup()

	memberlist, err := s.Tree.ListMembers(g.GetName())
	if err != nil {
		return nil, err
	}

	log.Printf("Membership of '%s' requested (%s@%s)",
		g.GetName(),
		client.GetService(),
		client.GetID())

	return &pb.EntityList{
		Members: memberlist,
	}, nil
}

func (s *NetAuthServer) ModifyGroupNesting(ctx context.Context, r *pb.ModGroupNestingRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	parent := r.GetParentGroup()
	child := r.GetChildGroup()
	mode := r.GetMode()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Either the entity must posses the right capability, or they
	// must be in the a group that is permitted to manage this one
	// based on membership.  Either is sufficient.
	if !s.manageByMembership(c.EntityID, child.GetName()) && !c.HasCapability("MODIFY_GROUP_MEMBERS") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.ModifyGroupExpansions(parent.GetName(), child.GetName(), mode); err != nil {
		return &pb.SimpleResult{Msg: proto.String("Membership could not be updated"), Success: proto.Bool(false)}, err
	}

	log.Printf("Group '%s'->'%s' expansion to '%s' by '%s' (%s@%s)",
		parent.GetName(),
		child.GetName(),
		mode,
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{Msg: proto.String("Nesting updated successfully"), Success: proto.Bool(true)}, nil
}
