package rpc

import (
	"context"

	pb "github.com/NetAuth/NetAuth/proto"
)

func (s *NetAuthServer) EntityInfo(ctx context.Context, netAuthRequest *pb.NetAuthRequest) (*pb.EntityMeta, error) {
	return &pb.EntityMeta{}, nil
}
