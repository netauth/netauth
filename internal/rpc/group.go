package rpc

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

// NewGroup creates a new group on the NetAuth server.  This action
// must be authorized by the presentation of a token containing
// appropriate capabilities.
func (s *NetAuthServer) NewGroup(ctx context.Context, r *pb.ModGroupRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	g := r.GetGroup()

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("CREATE_GROUP") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.NewGroup(g.GetName(), g.GetDisplayName(), g.GetManagedBy(), g.GetNumber()); err != nil {
		return nil, toWireError(err)
	}

	log.Printf("New Group '%s' created by '%s' (%s@%s)",
		g.GetName(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("New group created successfully"),
		Success: proto.Bool(true),
	}, toWireError(nil)
}

// DeleteGroup removes a group from the NetAuth server.  This action
// must be authorized by the presentation of a token containing
// apropriate capabilities.  This call will not CASCADE deletes and
// will not check if the group is empty before proceeding.  Other
// methods *should* safely handle this and check that they aren't
// pointing to a group that doesn't exist anymore, but its still good
// form to clean up references before calling this action.
func (s *NetAuthServer) DeleteGroup(ctx context.Context, r *pb.ModGroupRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	g := r.GetGroup()

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("DESTROY_GROUP") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.DeleteGroup(g.GetName()); err != nil {
		return nil, toWireError(err)
	}

	log.Printf("Group '%s' removed by '%s' (%s@%s)",
		g.GetName(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Group removed successfully"),
		Success: proto.Bool(true),
	}, toWireError(nil)
}

// GroupInfo returns as much information as is known about a group.
// This does not include group membership.
func (s *NetAuthServer) GroupInfo(ctx context.Context, r *pb.ModGroupRequest) (*pb.GroupInfoResult, error) {
	client := r.GetInfo()
	g := r.GetGroup()

	grp, err := s.Tree.GetGroupByName(g.GetName())
	if err != nil {
		return nil, toWireError(err)
	}

	log.Printf("Information on %s requested (%s@%s)",
		g.GetName(),
		client.GetService(),
		client.GetID())

	allGroups, err := s.Tree.ListGroups()
	if err != nil {
		log.Printf("Error summoning groups: %s", err)
	}
	var mgd []string
	for _, g := range allGroups {
		if g.GetManagedBy() == r.GetGroup().GetName() {
			mgd = append(mgd, g.GetName())
		}
	}

	return &pb.GroupInfoResult{Group: grp, Managed: mgd}, nil
}

// ModifyGroupMeta allows metadata stored on the group to be
// rewritten.  Some fields may not be changed using this action and
// must use more specialized calls which perform additional
// authorization and validation checks.  This action must be
// authorized by the presentation of a token containing appropriate
// capabilities.
func (s *NetAuthServer) ModifyGroupMeta(ctx context.Context, r *pb.ModGroupRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	g := r.GetGroup()

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Either the entity must posses the right capability, or they
	// must be in the a group that is permitted to manage this one
	// based on membership.  Either is sufficient.
	if !s.manageByMembership(c.EntityID, g.GetName()) && !c.HasCapability("MODIFY_GROUP_META") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.UpdateGroupMeta(g.GetName(), g); err != nil {
		return nil, toWireError(err)
	}

	log.Printf("Group '%s' modified by '%s' (%s@%s)",
		g.GetName(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Group modified successfully"),
		Success: proto.Bool(true),
	}, toWireError(err)
}
