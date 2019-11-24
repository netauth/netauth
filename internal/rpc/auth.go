package rpc

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"

	"github.com/netauth/netauth/internal/token"

	pb "github.com/netauth/protocol"
)

// AuthEntity performs entity authentication and returns boolean
// status for the authentication attempt.  This method should be
// preferred for systems that will not need a token, or will issue a
// token of their own on the authority of this response.
func (s *NetAuthServer) AuthEntity(ctx context.Context, r *pb.NetAuthRequest) (*pb.SimpleResult, error) {
	// Log out some useful stuff for debugging
	client := r.GetInfo()
	entity := r.GetEntity()
	s.Log.Info("Authenticating Entity",
		"entity", entity.GetID(),
		"service", client.GetService(),
		"client", client.GetID())

	// Construct and return the response.
	result := new(pb.SimpleResult)

	if err := s.Tree.ValidateSecret(entity.GetID(), entity.GetSecret()); err != nil {
		return nil, toWireError(err)
	}

	result.Success = proto.Bool(true)
	result.Msg = proto.String("Entity authentication succeeded")
	return result, toWireError(nil)
}

// GetToken is functionally identical to AuthEntity above, but will
// also return a token that can be used to perform further requests to
// the NetAuth server.
func (s *NetAuthServer) GetToken(ctx context.Context, r *pb.NetAuthRequest) (*pb.TokenResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()

	s.Log.Info("Token requested",
		"entity", e.GetID(),
		"service", client.GetService(),
		"client", client.GetID())

	// Run the normal authentication flow
	if err := s.Tree.ValidateSecret(e.GetID(), e.GetSecret()); err != nil {
		return nil, toWireError(err)
	}

	// Get the full fledged entity
	e, err := s.Tree.FetchEntity(e.GetID())
	if err != nil {
		s.Log.Warn("Entity vanished while fetching token!", "entity", e.GetID())
		return nil, toWireError(ErrInternalError)
	}

	// First get the capabilities that are provided by the entity
	// itself.
	caps := make(map[pb.Capability]int)
	if e.GetMeta() != nil {
		for _, c := range e.GetMeta().GetCapabilities() {
			caps[c]++
		}
	}

	// Next get the capabilities that are provided by any groups
	// the entity may be in; include indirects for authentication
	// queries.
	groupNames := s.Tree.GetMemberships(e, true)
	for _, name := range groupNames {
		g, err := s.Tree.FetchGroup(name)
		if err != nil {
			s.Log.Warn("Error loading group during GetToken", "group", name)
			continue
		}
		for _, c := range g.GetCapabilities() {
			caps[c]++
		}
	}

	// Flatten the capabilities out into a list
	var capabilities []pb.Capability
	for c := range caps {
		capabilities = append(capabilities, c)
	}

	// Successfully authenticated, now to construct a token
	claims := token.Claims{
		EntityID:     e.GetID(),
		Capabilities: capabilities,
	}

	// Generate the token with the specified claims
	tkn, err := s.Token.Generate(claims, token.GetConfig())
	if err != nil {
		return nil, toWireError(err)
	}

	// Construct the reply containing the token
	reply := pb.TokenResult{
		Success: proto.Bool(true),
		Msg:     proto.String("Token Granted"),
		Token:   &tkn,
	}

	return &reply, toWireError(nil)
}

// ValidateToken will attempt to determine the validity of a token
// previously issued by the NetAuth server.
func (s *NetAuthServer) ValidateToken(ctx context.Context, r *pb.NetAuthRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	s.Log.Info("Token validation requested",
		"entity", e.GetID(),
		"service", client.GetService(),
		"client", client.GetID())

	// Validate the token and if it validates, return that the
	// token is valid.
	_, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}

	return &pb.SimpleResult{
		Msg:     proto.String("Token verified"),
		Success: proto.Bool(true),
	}, toWireError(nil)
}

// ChangeSecret allows an entity secret to be reset.  There are two
// possible flows through this function based on whether or not the
// request is self-modifying or not.  In the case of a self modifying
// request (entity requests change of its own secret) then the entity
// must be in possession of the old secret, not a token, to authorize
// the change.  In the event the request is administrative (the entity
// is requesting the change of another entity's secret) then the
// entity must posses a token with the right capability.
func (s *NetAuthServer) ChangeSecret(ctx context.Context, r *pb.ModEntityRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	me := r.GetModEntity()
	t := r.GetAuthToken()

	if viper.GetBool("server.readonly") {
		s.Log.Warn("Denied ChangeSecret request (read-only mode is enabled)",
			"service", client.GetService(),
			"client", client.GetID())
		return &pb.SimpleResult{
			Msg:     proto.String("This server is in read-only mode"),
			Success: proto.Bool(false),
		}, toWireError(ErrReadOnly)
	}

	modself := false

	// Determine if this is a self modifying request or not
	if me != nil && e.GetID() == me.GetID() {
		modself = true
	}

	// Self modifying requests require the original password to
	// proceed
	if modself {
		err := s.Tree.ValidateSecret(e.GetID(), e.GetSecret())
		if err != nil {
			return nil, toWireError(err)
		}
		err = s.Tree.SetSecret(me.GetID(), me.GetSecret())
		if err != nil {
			return nil, toWireError(err)
		}
		s.Log.Info("Secret Changed",
			"entity", me.GetID(),
			"service", client.GetService(),
			"client", client.GetID())
		return &pb.SimpleResult{
			Msg:     proto.String("Secret Changed"),
			Success: proto.Bool(true),
		}, toWireError(nil)
	}

	// This change is being done administratively since modself
	// was false, so this needs a valid token to proceed.
	c, err := s.Token.Validate(t)
	if err != nil || !c.HasCapability(pb.Capability_CHANGE_ENTITY_SECRET) {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	// Change the secret per what was specified in the
	// modification entity struct.
	if err = s.Tree.SetSecret(me.GetID(), me.GetSecret()); err != nil {
		return nil, toWireError(err)
	}

	// Log this as an administrative change.
	s.Log.Info("Secret changed with authority",
		"entity", me.GetID(),
		"authority", c.EntityID,
		"service", client.GetService(),
		"client", client.GetID())
	return &pb.SimpleResult{
		Success: proto.Bool(true),
		Msg:     proto.String("Secret Changed"),
	}, toWireError(nil)
}

// ManageCapabilities permits the assignment and removal of
// capabilities from an entity or group.  If the entity and group are
// both specified, then the group will be ignored and the modification
// will be performed on the named entity.
func (s *NetAuthServer) ManageCapabilities(ctx context.Context, r *pb.ModCapabilityRequest) (*pb.SimpleResult, error) {
	entity := r.GetEntity()
	group := r.GetGroup()
	t := r.GetAuthToken()
	mode := r.GetMode()
	cap := r.GetCapability().String()

	if viper.GetBool("server.readonly") {
		s.Log.Warn("Denied ManageCapabilities request (read-only mode is enabled)",
			"service", r.GetInfo().GetService(),
			"client", r.GetInfo().GetID())
		return &pb.SimpleResult{
			Msg:     proto.String("This server is in read-only mode"),
			Success: proto.Bool(false),
		}, toWireError(ErrReadOnly)
	}

	// Validate the token and confirm the holder posses
	// GLOBAL_ROOT.  You might wonder why there isn't a capability
	// to assign other capabilities, but then you start going down
	// the rabbit hole and its much more straightforward to just
	// say that you need to be a global superuser to be able to
	// add more capabilities.
	c, err := s.Token.Validate(t)
	if err != nil || !c.HasCapability(pb.Capability_GLOBAL_ROOT) {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if entity != nil {
		switch mode {
		case "ADD":
			if err := s.Tree.SetEntityCapability(entity.GetID(), cap); err != nil {
				return &pb.SimpleResult{
					Success: proto.Bool(false),
					Msg:     proto.String("Error while adding capability"),
				}, toWireError(err)
			}
		case "DEL":
			if err := s.Tree.DropEntityCapability(entity.GetID(), cap); err != nil {
				return &pb.SimpleResult{
					Success: proto.Bool(false),
					Msg:     proto.String("Error while removing capability"),
				}, toWireError(err)
			}
		default:
			return &pb.SimpleResult{
				Success: proto.Bool(false),
				Msg:     proto.String("Mode must be either ADD or REMOVE"),
			}, toWireError(ErrMalformedRequest)
		}
	} else if group != nil {
		switch mode {
		case "ADD":
			if err := s.Tree.SetGroupCapability(group.GetName(), cap); err != nil {
				return &pb.SimpleResult{
					Success: proto.Bool(false),
					Msg:     proto.String("Error while adding capability"),
				}, toWireError(err)
			}
		case "DEL":
			if err := s.Tree.DropGroupCapability(group.GetName(), cap); err != nil {
				return &pb.SimpleResult{
					Success: proto.Bool(false),
					Msg:     proto.String("Error while removing capability"),
				}, toWireError(err)
			}
		default:
			return &pb.SimpleResult{
				Success: proto.Bool(false),
				Msg:     proto.String("Mode must be either ADD or REMOVE"),
			}, toWireError(ErrMalformedRequest)
		}

	} else {
		return &pb.SimpleResult{
			Success: proto.Bool(false),
			Msg:     proto.String("Either entity or group must be provided!"),
		}, toWireError(ErrMalformedRequest)
	}

	return &pb.SimpleResult{
		Success: proto.Bool(true),
		Msg:     proto.String("Capability Modified"),
	}, toWireError(nil)
}

// LockEntity locks an entity.  This action must be authorized with an
// appropriate token.
func (s *NetAuthServer) LockEntity(ctx context.Context, r *pb.NetAuthRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	if viper.GetBool("server.readonly") {
		s.Log.Warn("Denied LockEntity request (read-only mode is enabled)",
			"service", client.GetService(),
			"client", client.GetID())
		return &pb.SimpleResult{
			Msg:     proto.String("This server is in read-only mode"),
			Success: proto.Bool(false),
		}, toWireError(ErrReadOnly)
	}

	// Verify the correct capability is present in the token.
	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}
	if !c.HasCapability(pb.Capability_LOCK_ENTITY) {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.LockEntity(e.GetID()); err != nil {
		return &pb.SimpleResult{
			Success: proto.Bool(false),
			Msg:     proto.String("An error occurred while locking"),
		}, toWireError(err)
	}

	s.Log.Info("Entity Locked",
		"entity", e.GetID(),
		"authority", c.EntityID,
		"service", client.GetService(),
		"client", client.GetID())

	return &pb.SimpleResult{
		Success: proto.Bool(true),
		Msg:     proto.String("Entity is now locked"),
	}, toWireError(nil)
}

// UnlockEntity locks an entity.  This action must be authorized with an
// appropriate token.
func (s *NetAuthServer) UnlockEntity(ctx context.Context, r *pb.NetAuthRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	if viper.GetBool("server.readonly") {
		s.Log.Warn("Denied UnlockEntity request (read-only mode is enabled)",
			"service", client.GetService(),
			"client", client.GetID())
		return &pb.SimpleResult{
			Msg:     proto.String("This server is in read-only mode"),
			Success: proto.Bool(false),
		}, toWireError(ErrReadOnly)
	}

	// Verify the correct capability is present in the token.
	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}
	if !c.HasCapability(pb.Capability_UNLOCK_ENTITY) {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.UnlockEntity(e.GetID()); err != nil {
		return &pb.SimpleResult{
			Success: proto.Bool(false),
			Msg:     proto.String("An error occurred while locking"),
		}, toWireError(err)
	}

	s.Log.Info("Entity unlocked",
		"entity", e.GetID(),
		"authority", c.EntityID,
		"service", client.GetService(),
		"client", client.GetID())

	return &pb.SimpleResult{
		Success: proto.Bool(true),
		Msg:     proto.String("Entity is now unlocked"),
	}, toWireError(nil)
}
