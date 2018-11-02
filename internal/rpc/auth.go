package rpc

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"

	"github.com/NetAuth/NetAuth/internal/token"

	pb "github.com/NetAuth/Protocol"
)

// AuthEntity performs entity authentication and returns boolean
// status for the authentication attempt.  This method should be
// preferred for systems that will not need a token, or will issue a
// token of their own on the authority of this response.
func (s *NetAuthServer) AuthEntity(ctx context.Context, r *pb.NetAuthRequest) (*pb.SimpleResult, error) {
	// Log out some useful stuff for debugging
	client := r.GetInfo()
	entity := r.GetEntity()
	log.Printf("Authenticating %s (%s@%s)",
		entity.GetID(),
		client.GetService(),
		client.GetID())

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

	log.Printf("Token requested for %s (%s@%s)",
		e.GetID(),
		client.GetService(),
		client.GetID())

	// Run the normal authentication flow
	if err := s.Tree.ValidateSecret(e.GetID(), e.GetSecret()); err != nil {
		return nil, toWireError(err)
	}

	// Get the full fledged entity
	e, err := s.Tree.GetEntity(e.GetID())
	if err != nil {
		log.Println("Entity Vanished!")
		return nil, toWireError(ErrInternalError)
	}

	// First get the capabilities that are provided by the entity
	// itself.
	caps := make(map[string]int)
	if e.GetMeta() != nil {
		for _, c := range e.GetMeta().GetCapabilities() {
			caps[pb.Capability_name[int32(c)]]++
		}
	}

	// Next get the capabilities that are provided by any groups
	// the entity may be in; include indirects for authentication
	// queries.
	groupNames := s.Tree.GetMemberships(e, true)
	for _, name := range groupNames {
		g, err := s.Tree.GetGroupByName(name)
		if err != nil {
			log.Printf("Error loading group: %s", err)
			continue
		}
		for _, c := range g.GetCapabilities() {
			caps[pb.Capability_name[int32(c)]]++
		}
	}

	// Flatten the capabilities out into a list
	var capabilities []string
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

	log.Printf("Token validation requested by %s (%s@%s)",
		e.GetID(),
		client.GetService(),
		client.GetID())

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

	modself := false

	// Determine if this is a self modifying request or not
	if me != nil && e.GetID() == me.GetID() {
		modself = true
	}

	// Self modifying requests require the original password to
	// proceed
	err := s.Tree.ValidateSecret(e.GetID(), e.GetSecret())
	if modself && err == nil {
		err := s.Tree.SetEntitySecretByID(me.GetID(), me.GetSecret())
		if err != nil {
			return nil, toWireError(err)
		}
		log.Printf("Secret for %s changed (%s@%s)",
			me.GetID(),
			client.GetService(),
			client.GetID())
		return &pb.SimpleResult{
			Msg:     proto.String("Secret Changed"),
			Success: proto.Bool(true),
		}, toWireError(nil)
	}
	if err != nil {
		// Problem with the password in the self modifying
		// request, bail out!
		return nil, toWireError(err)
	}

	// This change is being done administratively since modself
	// was false, so this needs a valid token to proceed.
	c, err := s.Token.Validate(t)
	if err != nil || !c.HasCapability("CHANGE_ENTITY_SECRET") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	// Change the secret per what was specified in the
	// modification entity struct.
	if err = s.Tree.SetEntitySecretByID(me.GetID(), me.GetSecret()); err != nil {
		return nil, toWireError(err)
	}

	// Log this as an administrative change.
	log.Printf("Secret for %s administratively changed by %s (%s@%s)",
		me.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())
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

	// Validate the token and confirm the holder posses
	// GLOBAL_ROOT.  You might wonder why there isn't a capability
	// to assign other capabilities, but then you start going down
	// the rabbit hole and its much more straightforward to just
	// say that you need to be a global superuser to be able to
	// add more capabilities.
	c, err := s.Token.Validate(t)
	if err != nil || !c.HasCapability("GLOBAL_ROOT") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if entity != nil {
		switch mode {
		case "ADD":
			if err := s.Tree.SetEntityCapabilityByID(entity.GetID(), cap); err != nil {
				return &pb.SimpleResult{
					Success: proto.Bool(false),
					Msg:     proto.String("Error while adding capability"),
				}, toWireError(err)
			}
		case "REMOVE":
			if err := s.Tree.RemoveEntityCapabilityByID(entity.GetID(), cap); err != nil {
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
			if err := s.Tree.SetGroupCapabilityByName(group.GetName(), cap); err != nil {
				return &pb.SimpleResult{
					Success: proto.Bool(false),
					Msg:     proto.String("Error while adding capability"),
				}, toWireError(err)
			}
		case "REMOVE":
			if err := s.Tree.RemoveGroupCapabilityByName(group.GetName(), cap); err != nil {
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

	// Verify the correct capability is present in the token.
	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}
	if !c.HasCapability("LOCK_ENTITY") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.LockEntity(e.GetID()); err != nil {
		return &pb.SimpleResult{
			Success: proto.Bool(false),
			Msg: proto.String("An error occured while locking"),
		}, toWireError(err)
	}

	log.Printf("Entity %s locked by %s (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Success: proto.Bool(true),
		Msg: proto.String("Entity is now locked"),
	}, toWireError(nil)
}

// UnlockEntity locks an entity.  This action must be authorized with an
// appropriate token.
func (s *NetAuthServer) UnlockEntity(ctx context.Context, r *pb.NetAuthRequest) (*pb.SimpleResult, error) {
	client := r.GetInfo()
	e := r.GetEntity()
	t := r.GetAuthToken()

	// Verify the correct capability is present in the token.
	c, err := s.Token.Validate(t)
	if err != nil {
		return nil, toWireError(err)
	}
	if !c.HasCapability("LOCK_ENTITY") {
		return nil, toWireError(ErrRequestorUnqualified)
	}

	if err := s.Tree.UnlockEntity(e.GetID()); err != nil {
		return &pb.SimpleResult{
			Success: proto.Bool(false),
			Msg: proto.String("An error occured while locking"),
		}, toWireError(err)
	}

	log.Printf("Entity %s unlocked by %s (%s@%s)",
		e.GetID(),
		c.EntityID,
		client.GetService(),
		client.GetID())

	return &pb.SimpleResult{
		Success: proto.Bool(true),
		Msg: proto.String("Entity is now unlocked"),
	}, toWireError(nil)
}
