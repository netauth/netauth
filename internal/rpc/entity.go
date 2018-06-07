package rpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/NetAuth/NetAuth/internal/token"
	"github.com/NetAuth/NetAuth/pkg/errors"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

func (s *NetAuthServer) NewEntity(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("CREATE_ENTITY") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.NewEntity(e.GetID(), e.GetNumber(), e.GetSecret()); err != nil {
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

func (s *NetAuthServer) RemoveEntity(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("DESTROY_ENTITY") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.DeleteEntityByID(e.GetID()); err != nil {
		return &pb.SimpleResult{Success: proto.Bool(false), Msg: proto.String(fmt.Sprintf("%s", err))}, nil
	}

	log.Printf("Entity '%s' removed by '%s' (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Entity removed successfully"),
		Success: proto.Bool(true),
	}, nil
}

func (s *NetAuthServer) EntityInfo(ctx context.Context, r *pb.NetAuthRequest) (*pb.Entity, error) {
	client := r.GetInfo()
	e := r.GetEntity()

	log.Printf("Info requested on '%s' (%s@%s)",
		e.GetID(),
		client.GetService(),
		client.GetID())

	return s.Tree.GetEntity(e.GetID())
}

func (s *NetAuthServer) ModifyEntityMeta(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return &pb.SimpleResult{Msg: proto.String("Authentication Failure")}, nil
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("MODIFY_ENTITY_META") {
		return &pb.SimpleResult{Msg: proto.String("Requestor not qualified"), Success: proto.Bool(false)}, nil
	}

	if err := s.Tree.UpdateEntityMeta(e.GetID(), e.GetMeta()); err != nil {
		log.Printf("Metadata update error: %s", err)
		return nil, err
	}

	log.Printf("Metadata for '%s' by '%s' completed (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{Success: proto.Bool(true), Msg: proto.String("Metadata Updated")}, nil
}

func (s *NetAuthServer) ModifyEntityKeys(ctx context.Context, r *pb.ModEntityKeyRequest) (*pb.KeyList, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	mode := strings.ToUpper(r.GetMode())

	// If we aren't doing a read only operation then we need a
	// token for this
	var c token.Claims
	if mode != "LIST" {
		c, err := s.Token.Validate(t)
		if err != nil {
			return nil, err
		}

		// Verify the correct capability is present in the token or
		// that this is not a read only query.
		if !c.HasCapability("MODIFY_ENTITY_KEYS") {
			return nil, errors.E_NO_CAPABILITY
		}
	}

	// Get run the transaction on the key database.
	keys, err := s.Tree.UpdateEntityKeys(e.GetID(), r.GetMode(), r.GetType(), r.GetKey())
	if err != nil {
		return nil, err
	}

	verb := "updated"
	if mode == "LIST" {
		verb = "requested"
	}
	log.Printf("Keys for '%s' %s by '%s' (%s@%s)",
		e.GetID(),
		verb,
		c.EntityID,
		client.GetService(),
		client.GetID())
	return &pb.KeyList{Keys: keys}, nil
}
