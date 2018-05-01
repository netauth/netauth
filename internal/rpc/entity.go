package rpc

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/NetAuth/pkg/proto"
)

func (s *NetAuthServer) NewEntity(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// TODO(maldridge) check capabilities here so that you
	// actually need the right capability to do things.

	if err := s.Tree.NewEntity(e.GetID(), e.GetUidNumber(), e.GetSecret()); err != nil {
		return &pb.SimpleResult{Success: proto.Bool(false), Msg: proto.String(fmt.Sprintf("%s", err))}, nil
	}

	log.Printf("New entity '%s' created by '%s' (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("New entity created successfully"),
		Success: proto.Bool(true),
	}, nil
}
