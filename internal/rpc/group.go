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
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("CREATE_GROUP") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.NewGroup(g.GetName(), g.GetDisplayName(), g.GetManagedBy(), g.GetGidNumber()); err != nil {
		return &pb.SimpleResult{
			Msg:     proto.String("Group could not be created"),
			Success: proto.Bool(false),
		}, err
	}

	log.Printf("New Group '%s' created by '%s' (%s@%s)",
		g.GetName(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("New group created successfully"),
		Success: proto.Bool(true),
	}, nil
}

func (s *NetAuthServer) DeleteGroup(ctx context.Context, r *pb.ModGroupRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	g := r.GetGroup()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("DESTROY_GROUP") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.DeleteGroup(g.GetName()); err != nil {
		return &pb.SimpleResult{
			Msg:     proto.String("Group could not be removed"),
			Success: proto.Bool(false),
		}, err
	}

	log.Printf("Group '%s' removed by '%s' (%s@%s)",
		g.GetName(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Group removed successfully"),
		Success: proto.Bool(true),
	}, nil
}

func (s *NetAuthServer) ListGroups(ctx context.Context, r *pb.GroupListRequest) (*pb.GroupList, error) {
	client := r.GetInfo()

	list, err := s.Tree.ListGroups()
	if err != nil {
		return nil, err
	}

	log.Printf("Group list requested (%s@%s)",
		client.GetService(),
		client.GetID())

	return &pb.GroupList{Groups: list}, nil
}

func (s *NetAuthServer) ModifyGroupMeta(ctx context.Context, r *pb.ModGroupRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	t := r.GetAuthToken()
	g := r.GetGroup()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("MODIFY_GROUP_META") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.UpdateGroupMeta(g.GetName(), g); err != nil {
		return &pb.SimpleResult{
			Msg:     proto.String("Group could not be modified"),
			Success: proto.Bool(false),
		}, err
	}

	log.Printf("Group '%s' modified by '%s' (%s@%s)",
		g.GetName(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Group modified successfully"),
		Success: proto.Bool(true),
	}, nil
}
