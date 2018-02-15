package rpc

import (
	"context"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func (s *NetAuthServer) AddEntityToGroup(ctx context.Context, mer *pb.ModGroupDirectMembershipRequest) (*pb.SimpleResult, error) {

	if err := s.EM.AddEntityToGroup(mer); err != nil {
		return nil, err
	}
	return &pb.SimpleResult{Msg: proto.String("Membership updated"), Success: proto.Bool(true)}, nil
}

func (s *NetAuthServer) RemoveEntityFromGroup(ctx context.Context, mer *pb.ModGroupDirectMembershipRequest) (*pb.SimpleResult, error) {
	if err := s.EM.RemoveEntityFromGroup(mer); err != nil {
		return nil, err
	}
	return &pb.SimpleResult{Msg: proto.String("Membership updated"), Success: proto.Bool(true)}, nil
}
