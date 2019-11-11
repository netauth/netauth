package rpc

import (
	"context"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/netauth/netauth/internal/token"
	"github.com/spf13/viper"

	pb "github.com/NetAuth/Protocol"
)

// NewEntity creates a new entity.  This action must be authorized by
// the presentation of a valid token containing appropriate
// capabilities.
func (s *NetAuthServer) NewEntity(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	if viper.GetBool("server.readonly") {
		s.Log.Warn("Denied NewEntity request (read-only mode is enabled)",
			"service", client.GetService(),
			"client", client.GetID())
		return &pb.SimpleResult{
			Msg:     proto.String("This server is in read-only mode"),
			Success: proto.Bool(false),
		}, toWireError(ErrReadOnly)
	}

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability(pb.Capability_CREATE_ENTITY) {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.CreateEntity(e.GetID(), e.GetNumber(), e.GetSecret()); err != nil {
		return nil, toWireError(err)
	}

	s.Log.Info("New entity created",
		"entity", e.GetID(),
		"authority", c.EntityID,
		"service", client.GetService(),
		"client", client.GetID())

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

	if viper.GetBool("server.readonly") {
		s.Log.Warn("Denied RemoveEntity request (read-only mode is enabled)",
			"service", client.GetService(),
			"client", client.GetID())
		return &pb.SimpleResult{
			Msg:     proto.String("This server is in read-only mode"),
			Success: proto.Bool(false),
		}, toWireError(ErrReadOnly)
	}

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability(pb.Capability_DESTROY_ENTITY) {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.DestroyEntity(e.GetID()); err != nil {
		return nil, toWireError(err)
	}

	s.Log.Info("Entity removed",
		"entity", e.GetID(),
		"authority", c.EntityID,
		"service", client.GetService(),
		"client", client.GetID())

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

	s.Log.Info("Entity information request",
		"entity", e.GetID(),
		"service", client.GetService(),
		"client", client.GetID())

	e, err := s.Tree.FetchEntity(e.GetID())
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

	if viper.GetBool("server.readonly") {
		s.Log.Warn("Denied ModifyEntityMeta request (read-only mode is enabled)",
			"service", client.GetService(),
			"client", client.GetID())
		return &pb.SimpleResult{
			Msg:     proto.String("This server is in read-only mode"),
			Success: proto.Bool(false),
		}, toWireError(ErrReadOnly)
	}

	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	// Verify the correct capability is present in the token.
	if !c.HasCapability(pb.Capability_MODIFY_ENTITY_META) {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.UpdateEntityMeta(e.GetID(), e.GetMeta()); err != nil {
		s.Log.Error("Metadata update error", "error", err)
		return nil, toWireError(err)
	}

	s.Log.Info("Entity metadata update",
		"entity", e.GetID(),
		"authority", c.EntityID,
		"service", client.GetService(),
		"client", client.GetID())

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

	if viper.GetBool("server.readonly") && mode != "LIST" {
		s.Log.Warn("Denied ModifyEntityKeys request (read-only mode is enabled)",
			"service", client.GetService(),
			"client", client.GetID())
		return &pb.KeyList{}, toWireError(ErrReadOnly)
	}

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
		if !c.HasCapability(pb.Capability_MODIFY_ENTITY_KEYS) && c.EntityID != e.GetID() {
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
	s.Log.Info("Entity Key Handling Event",
		"entity", e.GetID(),
		"authority", c.EntityID,
		"action", verb,
		"service", client.GetService(),
		"client", client.GetID())
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

	if viper.GetBool("server.readonly") && mode != "READ" {
		s.Log.Warn("Denied ModifyUntypedEntityMeta request (read-only mode is enabled)",
			"service", client.GetService(),
			"client", client.GetID())
		return &pb.UntypedMetaResult{}, toWireError(ErrReadOnly)
	}

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
		if !c.HasCapability(pb.Capability_MODIFY_ENTITY_META) {
			return nil, toWireError(ErrRequestorUnqualified)
		}
	}

	meta, err := s.Tree.ManageUntypedEntityMeta(e.GetID(), r.GetMode(), r.GetKey(), r.GetValue())
	if err != nil {
		return nil, toWireError(err)
	}

	if mode != "READ" {
		s.Log.Info("Entity UntypedMeta updated",
			"entity", e.GetID(),
			"authority", c.EntityID,
			"service", client.GetService(),
			"client", client.GetID(),
		)
	}

	return &pb.UntypedMetaResult{UntypedMeta: meta}, toWireError(nil)
}
