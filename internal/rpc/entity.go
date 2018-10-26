package rpc

import (
	"context"
	"log"
	"strings"

	"github.com/NetAuth/NetAuth/internal/token"
	"github.com/golang/protobuf/proto"

	pb "github.com/NetAuth/Protocol"
)

// NewEntity creates a new entity.  This action must be authorized by
// the presentation of a valid token containing appropriate
// capabilities.
func (s *NetAuthServer) NewEntity(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("CREATE_ENTITY") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.NewEntity(e.GetID(), e.GetNumber(), e.GetSecret()); err != nil {
		return nil, toWireError(err)
	}

	log.Printf("New entity '%s' created by '%s' (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("New entity created successfully"),
		Success: proto.Bool(true),
	}, toWireError(nil)
}

// RemoveEntity removes an entity.  This action must be authorized by
// the presentation of a valid token containing appropriate
// capabilities.
func (s *NetAuthServer) RemoveEntity(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("DESTROY_ENTITY") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.DeleteEntityByID(e.GetID()); err != nil {
		return nil, toWireError(err)
	}

	log.Printf("Entity '%s' removed by '%s' (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Entity removed successfully"),
		Success: proto.Bool(true),
	}, toWireError(nil)
}

// EntityInfo returns as much information about an entity is as known.
// This response will not include information about the entity's
// memberships in groups within the tree, but will include all fields
// in the EntityMeta section.
func (s *NetAuthServer) EntityInfo(ctx context.Context, r *pb.NetAuthRequest) (*pb.Entity, error) {
	client := r.GetInfo()
	e := r.GetEntity()

	log.Printf("Info requested on '%s' (%s@%s)",
		e.GetID(),
		client.GetService(),
		client.GetID())

	e, err := s.Tree.GetEntity(e.GetID())
	return e, toWireError(err)
}

// ModifyEntityMeta can be used to modify the EntityMeta section of an
// Entity.  This request must be authorized by a token that contains
// the correct capabilities to modify others.  Some fields cannot be
// changed by this mechanism and must be changed via other calls which
// perform more authorization and validation checks.
func (s *NetAuthServer) ModifyEntityMeta(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability("MODIFY_ENTITY_META") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.UpdateEntityMeta(e.GetID(), e.GetMeta()); err != nil {
		log.Printf("Metadata update error: %s", err)
		return nil, toWireError(err)
	}

	log.Printf("Metadata for '%s' by '%s' completed (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Msg:     proto.String("Metadata Updated"),
		Success: proto.Bool(true),
	}, toWireError(nil)
}

// ModifyEntityKeys can be used to add, remove, or retrieve the keys
// associated with an entity.  This action must be authorized by the
// presentation of a token with appropriate capabilities.
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
			return nil, toWireError(err)
		}

		// Verify the correct capability is present in the token or
		// that this is not a read only query.
		if !c.HasCapability("MODIFY_ENTITY_KEYS") {
			return nil, toWireError(ErrRequestorUnqualified)
		}
	}

	// Get run the transaction on the key database.
	keys, err := s.Tree.UpdateEntityKeys(e.GetID(), r.GetMode(), r.GetType(), r.GetKey())
	if err != nil {
		return nil, toWireError(err)
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
	return &pb.KeyList{
		Keys: keys,
	}, toWireError(nil)
}

// ModifyUntypedEntityMeta alters the data stored in the untyped K/V
// segment of an entity's metadata.  This action must be authorized by
// the presentation of a token with appropriate capabilities.
func (s *NetAuthServer) ModifyUntypedEntityMeta(ctx context.Context, r *pb.ModEntityMetaRequest) (*pb.UntypedMetaResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	mode := strings.ToUpper(r.GetMode())

	// If we aren't doing a read only operation then we need a
	// token for this
	var c token.Claims
	var err error
	if mode != "READ" {
		c, err = s.Token.Validate(t)
		if err != nil {
			return nil, toWireError(err)
		}

		// Verify the correct capability is present in the token or
		// that this is not a read only query.
		if !c.HasCapability("MODIFY_ENTITY_KEYS") {
			return nil, toWireError(ErrRequestorUnqualified)
		}
	}

	meta, err := s.Tree.ManageUntypedEntityMeta(e.GetID(), r.GetMode(), r.GetKey(), r.GetValue())
	if err != nil {
		return nil, toWireError(err)
	}

	if mode != "READ" {
		log.Printf("UntypedMeta for %s updated by %s (%s@%s)",
			e.GetID(),
			c.EntityID,
			client.GetService(),
			client.GetID(),
		)
	}

	return &pb.UntypedMetaResult{UntypedMeta: meta}, toWireError(nil)
}
