package rpc

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

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
		return nil, toWireError(RequestorUnqualified)
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
		return nil, toWireError(RequestorUnqualified)
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
		return nil, toWireError(RequestorUnqualified)
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
